package gateway

import (
	"errors"

	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
)

var errIncorrectResponseType = errors.New("incorrect response type")

type Bridge interface {
	Add(*cluster.Peer, pb.Stoa_KeepServer) error
}

type server struct {
	pb.UnimplementedStoaServer
	logger logging.Logger
	peer   *cluster.Peer
	bridge Bridge
}

func newServer(p *cluster.Peer, b Bridge) *server {
	return &server{
		logger: logging.NewLogger("stoa"),
		peer:   p,
		bridge: b,
	}
}

func (s *server) Keep(stream pb.Stoa_KeepServer) error {
	return s.bridge.Add(s.peer, stream)
}
