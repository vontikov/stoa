package cluster

import (
	"github.com/google/uuid"
	"github.com/vontikov/stoa/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (f *FSM) Add(peer *Peer, stream pb.Stoa_KeepServer) error {
	id := uuid.New().String()
	f.streams.Put(id, stream)
	f.logger.Debug("Stream added", "id", id)

	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			f.streams.Remove(id)
			f.logger.Debug("Stream removed", "id", id)
			return ctx.Err()
		default:
			p, err := stream.Recv()
			if err != nil {
				if s, ok := status.FromError(err); ok {
					if s.Code() == codes.Canceled {
						return nil
					}
				}
				f.logger.Error("Stream error", "message", err)
				return err
			}
			if peer.IsLeader() {
				f.ping(p)
			}
		}
	}
}

func (f *FSM) ping(p *pb.Ping) {
	f.logger.Debug("Ping received", "message", p)
	keys := f.ms.Keys()
	for _, k := range keys {
		v := f.ms.Get(k)
		if v == nil {
			continue
		}
		mx := v.(*mutexRecord)
		mx.touch(p.Id)
	}
}
