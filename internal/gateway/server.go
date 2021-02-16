package gateway

import (
	"bytes"
	"context"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
)

var checkLeadershipPeriod = 500 * time.Millisecond

type streamRecord struct {
	id []byte
	ch chan *pb.Status
}

type server struct {
	pb.UnimplementedStoaServer
	logger  logging.Logger
	ctx     context.Context
	cluster cluster.Cluster

	mu      sync.RWMutex // protects streams
	streams []streamRecord
}

func newServer(ctx context.Context, c cluster.Cluster, logger logging.Logger) *server {
	s := &server{
		logger:  logger,
		ctx:     ctx,
		cluster: c,
	}
	s.startStatusPoller()
	return s
}

func (s *server) Watch(v *pb.ClientId, stream pb.Stoa_WatchServer) error {
	if err := v.Validate(); err != nil {
		return err
	}

	if !s.cluster.IsLeader() {
		return ErrNotLeader
	}

	if err := s.handshake(v.Id, stream); err != nil {
		return err
	}

	id := v.Id
	ch := make(chan *pb.Status, 1) // needs some room in case of lost leadership
	s.mu.Lock()
	s.streams = append(s.streams, streamRecord{id, ch})
	s.mu.Unlock()
	s.logger.Debug("watcher arrived", "id", string(id))

	defer func() {
		s.mu.Lock()
		for i, r := range s.streams {
			if bytes.Equal(r.id, id) {
				s.streams = append(s.streams[:i], s.streams[i+1:]...)
				break
			}
		}
		s.mu.Unlock()
		s.logger.Debug("watcher gone", "id", string(id))
	}()

	checkLeadership := time.NewTicker(checkLeadershipPeriod)
	defer checkLeadership.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-checkLeadership.C:
			if !s.cluster.IsLeader() {
				ch <- &pb.Status{U: &pb.Status_C{C: &pb.ClusterStatus{S: pb.ClusterStatus_LEADERSHIP_LOST}}}
			}
		case st := <-ch:
			if !s.cluster.IsLeader() {
				if s.logger.IsTrace() {
					s.logger.Trace("leadership lost", "id", id)
				}
				return ErrNotLeader
			}
			if err := stream.Send(st); err != nil {
				s.logger.Warn("status error", "message", err)
				break
			}
			if s.logger.IsTrace() {
				s.logger.Trace("status sent", "id", string(id), "message", st)
			}
		}
	}
}

func (s *server) handshake(clientID []byte, stream pb.Stoa_WatchServer) (err error) {
	if err = stream.Send(&pb.Status{U: &pb.Status_Id{Id: &pb.ClientId{Id: clientID}}}); err != nil {
		s.logger.Debug("watch handshake failed", "id", clientID, "message", err)
		return
	}
	if s.logger.IsTrace() {
		s.logger.Trace("watch handshake", "id", string(clientID))
	}
	return
}

func (s *server) startStatusPoller() {
	ch := s.cluster.Status()
	go func() {
		defer s.logger.Debug("status poller stopped")
		for {
			select {
			case <-s.ctx.Done():
				return
			case status := <-ch:
				s.mu.RLock()
				for _, r := range s.streams {
					select {
					case r.ch <- status:
					default:
						if s.logger.IsTrace() {
							s.logger.Trace("watch channel blocked")
						}
					}
				}
				s.mu.RUnlock()
			}
		}
	}()
	s.logger.Debug("status poller started")
}

func (s *server) Ping(ctx context.Context, v *pb.ClientId) (*pb.Empty, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	msg := pb.ClusterCommand{
		Command: pb.ClusterCommand_SERVICE_PING,
		Payload: &pb.ClusterCommand_ClientId{ClientId: v},
	}

	cmd, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	fut := s.cluster.Raft().Apply(cmd, 0)
	if err := fut.Error(); err != nil {
		return nil, wrapError(err)
	}
	return &pb.Empty{}, nil
}
