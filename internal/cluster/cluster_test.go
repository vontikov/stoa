package cluster

import (
	"testing"

	"context"
	"strings"
	"sync"
	"time"

	"github.com/stretchr/testify/assert"
	cc "github.com/vontikov/go-concurrent"
	"github.com/vontikov/stoa/internal/logging"
)

type R struct {
	cluster Cluster
	cancel  context.CancelFunc
}

func TestCluster(t *testing.T) {

	assert := assert.New(t)
	logging.SetLevel("trace")

	const (
		addr1 = "127.0.0.1:3501"
		addr2 = "127.0.0.1:3502"
		addr3 = "127.0.0.1:3503"
		addr4 = "127.0.0.1:3504"
		addr5 = "127.0.0.1:3505"
	)

	peers := []string{
		strings.Join([]string{addr1, addr2, addr3, addr4, addr5}, PeerListSep),
		strings.Join([]string{addr2}, PeerListSep),
		strings.Join([]string{addr3}, PeerListSep),
		strings.Join([]string{addr4}, PeerListSep),
		strings.Join([]string{addr5}, PeerListSep),
	}

	clusters := cc.NewSynchronizedMap(0)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		clusters.Range(func(k, v interface{}) bool {
			return assert.Nil(<-v.(R).cluster.Done())
		})
	}()

	var wg sync.WaitGroup
	for i, p := range peers {
		wg.Add(1)
		go func(p string, i int) {
			cctx, ccancel := context.WithCancel(ctx)
			cluster, err := New(cctx, WithPeers(p))
			assert.Nil(err)
			clusters.Put(i, R{cluster, ccancel})
			wg.Done()
		}(p, i)
	}
	wg.Wait()

	ensureSameLeader(t, clusters)
	assert.True(clusters.Get(0).(R).cluster.IsLeader())

	// drain signals
	clusters.Range(func(k, v interface{}) bool {
		ch := v.(R).cluster.WatchLeadership()
		for len(ch) > 0 {
			<-ch
		}
		return true
	})

	// transfer leadership
	assert.Nil(clusters.Get(0).(R).cluster.LeadershipTransfer())

	// make sure there is a new leader
	timeout := time.After(10 * time.Second)
loop:
	for {
		select {
		case <-timeout:
			t.Logf("election timeout")
			t.Fail()
		default:
			for i := 1; i < clusters.Size(); i++ {
				ch := clusters.Get(i).(R).cluster.WatchLeadership()
				if len(ch) > 0 && <-ch {
					break loop
				}
			}
		}
	}
	ensureSameLeader(t, clusters)
}

func TestClusterLeaderRestart(t *testing.T) {
	assert := assert.New(t)
	logging.SetLevel("info")

	const (
		addr1 = "127.0.0.1:3501"
		addr2 = "127.0.0.1:3502"
		addr3 = "127.0.0.1:3503"
	)

	peers := []string{
		strings.Join([]string{addr1, addr2, addr3}, PeerListSep),
		strings.Join([]string{addr2}, PeerListSep),
		strings.Join([]string{addr3}, PeerListSep),
	}

	clusters := cc.NewSynchronizedMap(0)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		clusters.Range(func(k, v interface{}) bool {
			return assert.Nil(<-v.(R).cluster.Done())
		})
	}()

	var wg sync.WaitGroup
	for i, p := range peers {
		wg.Add(1)
		go func(p string, i int) {
			cctx, ccancel := context.WithCancel(ctx)
			cluster, err := New(cctx, WithPeers(p))
			assert.Nil(err)
			clusters.Put(i, R{cluster, ccancel})
			wg.Done()
		}(p, i)
	}
	wg.Wait()
	ensureSameLeader(t, clusters)

	r := clusters.Get(0).(R)
	r.cancel()
	assert.Nil(<-r.cluster.Done())

	cctx, ccancel := context.WithCancel(ctx)
	cluster, err := New(cctx, WithPeers(peers[0]))
	assert.Nil(err)
	clusters.Put(0, R{cluster, ccancel})

	ensureSameLeader(t, clusters)
}

func ensureSameLeader(t *testing.T, clusters cc.Map) {
	timeout := time.After(10 * time.Second)
	expectedLeader := ""
	for {
		select {
		case <-timeout:
			t.Logf("election timeout")
			t.Fail()
		default:
			r := true
			clusters.Range(func(k, v interface{}) bool {
				c := v.(R).cluster
				if c.LeaderAddress() != expectedLeader {
					expectedLeader = c.LeaderAddress()
					r = false
					return false
				}
				return true
			})
			if r {
				return
			}
		}
	}
}
