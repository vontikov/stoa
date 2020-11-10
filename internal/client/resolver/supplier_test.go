package resolver

import (
	"testing"

	"context"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"github.com/vontikov/stoa/internal/discovery"
	"github.com/vontikov/stoa/pkg/pb"
)

func TestSupplier(t *testing.T) {
	const (
		peerID = "id"

		grpcIP         = "127.0.0.1"
		grpcPort int32 = 3501

		discoveryIP   = "224.0.0.1"
		discoveryPort = 1237
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ms := func() ([]byte, error) {
		return proto.Marshal(
			&pb.Discovery{
				Id:       peerID,
				GrpcIp:   grpcIP,
				GrpcPort: grpcPort,
			})
	}
	sender, err := discovery.NewSender(discoveryIP, discoveryPort, ms)
	assert.Nil(t, err)

	suppl, err := NewDiscoverySupplier(discoveryIP, discoveryPort)
	assert.Nil(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sender.Run(ctx)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		suppl.Run(ctx)
	}()

	go func() {
		for {
			ch := suppl.Get("")
			for p := range ch {
				assert.Equal(t, grpcIP, p.GrpcIP)
				assert.Equal(t, grpcPort, p.GrpcPort)
				cancel()
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	wg.Wait()
}
