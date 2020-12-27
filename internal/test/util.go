package test

import (
	"testing"

	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/gateway"
	"github.com/vontikov/stoa/internal/logging"
)

const (
	testLogLevel = "debug"
)

func init() {
	logging.SetLevel(testLogLevel)
}

func RunTestCluster(ctx context.Context, t *testing.T, basePort int) []cluster.Cluster {
	bindAddr1 := fmt.Sprintf("127.0.0.1:%d", basePort+21)
	bindAddr2 := fmt.Sprintf("127.0.0.1:%d", basePort+22)
	bindAddr3 := fmt.Sprintf("127.0.0.1:%d", basePort+23)

	peers := []struct {
		grpcPort int
		httpPort int
		peers    string
	}{

		{basePort + 1, basePort + 11, strings.Join([]string{bindAddr1, bindAddr2, bindAddr3}, cluster.PeerListSep)},
		{basePort + 2, basePort + 12, strings.Join([]string{bindAddr2}, cluster.PeerListSep)},
		{basePort + 3, basePort + 14, strings.Join([]string{bindAddr3}, cluster.PeerListSep)},
	}

	var wg sync.WaitGroup
	var m sync.Mutex
	var clusters []cluster.Cluster
	for _, p := range peers {
		wg.Add(1)
		go func(grpcPort, httpPort int, peers string) {
			cluster, err := cluster.New(cluster.WithPeers(peers))
			if err != nil {
				t.Logf("cluster error: %v", err)
				t.Fail()
			}
			m.Lock()
			clusters = append(clusters, cluster)
			m.Unlock()
			gateway, err := gateway.New(ctx, cluster,
				gateway.WithListenAddress("localhost"),
				gateway.WithGRPCPort(grpcPort),
				gateway.WithHTTPPort(httpPort))
			if err != nil {
				t.Logf("gateway error: %v", err)
				t.Fail()
			}
			wg.Done()

			<-ctx.Done()
			gateway.Wait()
			cluster.Shutdown()
		}(p.grpcPort, p.httpPort, p.peers)
	}
	wg.Wait()

	// make sure the leader is the same on all the peers
	timeout := time.After(10 * time.Second)
loop:
	for {
		select {
		case <-timeout:
			t.Logf("Election timeout")
			t.Fail()
		default:
			for _, c := range clusters {
				if c.LeaderAddress() != bindAddr1 {
					continue loop
				}
			}
			break loop
		}
	}
	return clusters
}
