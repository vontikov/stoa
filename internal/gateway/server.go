package gateway

import (
	"errors"

	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrIncorrectResponseType = errors.New("incorrect response type")
var ErrNotLeader = errors.New("not leader")

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

func (s *server) Keep(stream pb.Stoa_KeepServer) error {
	if !s.cluster.IsLeader() {
		return status.Errorf(codes.FailedPrecondition, ErrNotLeader.Error())
	}

	return s.cluster.AddKeeper(stream)
}
