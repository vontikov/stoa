package gateway

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
)

var checkLeadershipPeriod = 500 * time.Millisecond

type server struct {
	pb.UnimplementedStoaServer
	logger  logging.Logger
	ctx     context.Context
	cluster cluster.Cluster
}

func newServer(ctx context.Context, c cluster.Cluster, logger logging.Logger) *server {
	s := &server{
		logger:  logger,
		ctx:     ctx,
		cluster: c,
	}
	return s
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
