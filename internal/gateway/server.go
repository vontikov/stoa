package gateway

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
)

type server struct {
	pb.UnimplementedStoaServer
	logger  logging.Logger
	cluster cluster.Cluster
}

func newServer(p cluster.Cluster) *server {
	return &server{
		logger:  logging.NewLogger("stoa"),
		cluster: p,
	}
}

func (s *server) Ping(ctx context.Context, v *pb.ClientId) (*pb.Empty, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	msg := pb.ClusterCommand{
		Command: pb.ClusterCommand_SERVICE_PING,
		Payload: &pb.ClusterCommand_Cid{Cid: v},
	}

	if err := processMetadata(ctx, &msg); err != nil {
		return nil, err
	}
	if msg.BarrierEnabled {
		fut := s.cluster.Raft().Barrier(time.Duration(msg.BarrierMillis) * time.Millisecond)
		if err := fut.Error(); err != nil {
			return nil, err
		}
	}

	cmd, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	fut := s.cluster.Raft().Apply(cmd, 0)
	if err := fut.Error(); err != nil {
		return nil, wrapError(err)
	}
	if r, ok := fut.Response().(*pb.Empty); ok {
		return r, nil
	}
	panic(ErrIncorrectResponseType)
}
