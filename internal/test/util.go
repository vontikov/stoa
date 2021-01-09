package test

import (
	"errors"

	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

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

func StartCluster(ctx context.Context, basePort int, peerNum int) (peers []cluster.Cluster, bootstrap string, err error) {
	if peerNum <= 0 {
		panic("peer number must be greater then 0")
	}
	if peerNum > 9 {
		panic("peer number must not be greater then 9")
	}
	if peerNum%2 == 0 {
		panic("peer number must not be even")
	}

	const (
		grpcPortShift = 1
		httpPortShift = 21
		bindPortShift = 11
	)

	type peer struct {
		grpcPort int
		httpPort int
		peers    string
	}

	var bindAddrs, bsAddrs []string
	var peerList []peer

	for i := 0; i < peerNum; i++ {
		bsAddrs = append(bsAddrs, fmt.Sprintf("127.0.0.1:%d", basePort+grpcPortShift+i))
		bindAddrs = append(bindAddrs, fmt.Sprintf("127.0.0.1:%d", basePort+bindPortShift+i))
	}

	for i := 0; i < peerNum; i++ {
		b := bindAddrs[i]
		if i == 0 {
			b = strings.Join(bindAddrs, cluster.PeerListSep)
		}
		peerList = append(peerList, peer{basePort + grpcPortShift + i, basePort + i + httpPortShift + i, b})
	}

	var m sync.Mutex
	var g errgroup.Group

	for i, p := range peerList {
		idx := i
		p := p
		g.Go(func() error {
			cluster, err := cluster.New(ctx,
				cluster.WithPeers(p.peers),
				cluster.WithLoggerName(fmt.Sprintf("raft-%d", idx)),
			)
			if err != nil {
				return err
			}
			m.Lock()
			peers = append(peers, cluster)
			m.Unlock()
			gateway, err := gateway.New(ctx, cluster,
				gateway.WithListenAddress("localhost"),
				gateway.WithGRPCPort(p.grpcPort),
				gateway.WithHTTPPort(p.httpPort))
			if err != nil {
				return err
			}

			go func() {
				<-ctx.Done()
				gateway.Wait()
			}()

			return nil
		})
	}

	err = g.Wait()
	if err != nil {
		return
	}

	// make sure the leader is the same on all the peers
	m.Lock()
	defer m.Unlock()
	leader, err := GetLeader(peers, 10*time.Second)
	if err != nil {
		return
	}
	if leader != bindAddrs[0] {
		err = errors.New("unexpected leader address")
	}

	bootstrap = strings.Join(bsAddrs, cluster.PeerListSep)
	return
}

func GetLeader(clusters []cluster.Cluster, timeout time.Duration) (leader string, err error) {
	t := time.After(timeout)
	for {
		select {
		case <-t:
			err = errors.New("election timeout")
			return
		default:
			r := true
			for _, c := range clusters {
				l := c.LeaderAddress()
				if l == "" {
					r = false
					break
				}
				if l != leader {
					leader = l
					r = false
					break
				}
			}
			if r {
				return
			}
		}
	}
}
