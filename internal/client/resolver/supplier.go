package resolver

import (
	"context"
	"net"
	"sync"

	"github.com/golang/protobuf/proto"

	"github.com/vontikov/stoa/internal/common"
	"github.com/vontikov/stoa/internal/discovery"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
)

type peer struct {
	ID       string
	GrpcIP   string
	GrpcPort int32
	Leader   bool
}

type Supplier interface {
	common.Runnable
	Get(string) <-chan peer
}

type discoverySupplier struct {
	logger   logging.Logger
	peers    *sync.Map
	receiver common.Runnable
}

func NewDiscoverySupplier(ip string, port int) (Supplier, error) {
	peers := sync.Map{}
	receiver, err := discovery.NewReceiver(ip, port, handler(&peers))
	if err != nil {
		return nil, err
	}

	return &discoverySupplier{
		logger:   logging.NewLogger("supplier"),
		peers:    &peers,
		receiver: receiver,
	}, nil
}

func (s *discoverySupplier) Run(ctx context.Context) error {
	return s.receiver.Run(ctx)
}

func (s *discoverySupplier) Get(endpoint string) <-chan peer {
	ch := make(chan peer)
	go func() {
		s.peers.Range(func(k, v interface{}) bool {
			if p, ok := v.(peer); ok {
				ch <- p
				s.logger.Debug("peer", "data", p)
				return true
			}
			panic("must be peer instance")
		})
		close(ch)
	}()
	return ch
}

func handler(peers *sync.Map) discovery.Handler {
	dm := pb.Discovery{}
	return func(src *net.UDPAddr, msg []byte) error {
		if err := proto.Unmarshal(msg, &dm); err != nil {
			return err
		}
		peer := peer{
			ID:       dm.Id,
			GrpcIP:   dm.GrpcIp,
			GrpcPort: dm.GrpcPort,
			Leader:   dm.Leader,
		}
		peers.Store(dm.Id, peer)
		return nil
	}
}
