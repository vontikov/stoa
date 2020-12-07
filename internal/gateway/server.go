package gateway

import (
	"errors"

	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
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
