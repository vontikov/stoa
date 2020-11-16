package gateway

import (
	"errors"

	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
)

var errIncorrectResponseType = errors.New("incorrect response type")

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
	return s.cluster.AddKeeper(stream)
}
