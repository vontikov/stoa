package client

import (
	"testing"

	"context"
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

func runTestCluster(ctx context.Context, t *testing.T) {
	const (
		addr1 = "127.0.0.1:3501"
		addr2 = "127.0.0.1:3502"
		addr3 = "127.0.0.1:3503"
	)

	peers := []struct {
		grpcPort int
		httpPort int
		peers    string
	}{

		{2501, 2601, strings.Join([]string{addr1, addr2, addr3}, cluster.PeerListSep)},
		{2502, 2602, strings.Join([]string{addr2}, cluster.PeerListSep)},
		{2503, 2603, strings.Join([]string{addr3}, cluster.PeerListSep)},
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
			gateway, err := gateway.New("localhost", grpcPort, httpPort, cluster)
			if err != nil {
				t.Logf("gateway error: %v", err)
				t.Fail()
			}
			wg.Done()

			<-ctx.Done()
			gateway.Shutdown()
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
				if c.LeaderAddress() != addr1 {
					continue loop
				}
			}
			break loop
		}
	}
}

func TestClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	runTestCluster(ctx, t)
}
