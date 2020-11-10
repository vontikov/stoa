package discovery

import (
	"net"

	"google.golang.org/protobuf/proto"

	"github.com/vontikov/go-concurrent"
	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/pkg/pb"
)

// Hook defines a function called on any event
type Hook func(...interface{})

// DefaultHandler returns default Handler
func DefaultHandler(p *cluster.Peer, onDiscoveryHook Hook) Handler {
	peers := concurrent.NewSynchronizedMap(0)

	return func(src *net.UDPAddr, msg []byte) error {
		dm := pb.Discovery{}
		if err := proto.Unmarshal(msg, &dm); err != nil {
			return err
		}

		if dm.Id == p.ID() || !p.IsLeader() {
			return nil
		}

		if ok := peers.PutIfAbsent(dm.Id, src); ok {
			if onDiscoveryHook != nil {
				onDiscoveryHook(src, &dm)
			}
			if err := p.AddVoter(dm.Id, dm.BindIp, int(dm.BindPort)); err != nil {
				return err
			}
		}
		return nil
	}
}

// DefaultMessageSupplier returns default message supplier
func DefaultMessageSupplier(p *cluster.Peer, grpcIP string, grpcPort int) MessageSupplier {
	return func() ([]byte, error) {
		msg, err := proto.Marshal(
			&pb.Discovery{
				Id:       p.ID(),
				GrpcIp:   grpcIP,
				GrpcPort: int32(grpcPort),
				BindIp:   p.BindIP(),
				BindPort: int32(p.BindPort()),
				Leader:   p.IsLeader(),
			})
		if err != nil {
			return nil, err
		}
		return msg, nil
	}
}
