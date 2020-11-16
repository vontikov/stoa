package cluster

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vontikov/stoa/internal/logging"
)

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
		strings.Join([]string{addr1, addr2, addr3, addr4, addr5}, peerListSep),
		strings.Join([]string{addr2}, peerListSep),
		strings.Join([]string{addr3}, peerListSep),
		strings.Join([]string{addr4}, peerListSep),
		strings.Join([]string{addr5}, peerListSep),
	}

	var wg sync.WaitGroup
	var m sync.Mutex
	var clusters []Cluster

	for _, p := range peers {
		wg.Add(1)
		go func(p string) {
			cluster, err := New(WithPeers(p))
			assert.Nil(err)
			m.Lock()
			clusters = append(clusters, cluster)
			m.Unlock()
			wg.Done()
		}(p)
	}
	wg.Wait()
	m.Lock()
	defer m.Unlock()

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

	// transfer the leadership
	for _, c := range clusters {
		if c.IsLeader() {
			// drain the previous signals
			ch := c.WatchLeadership()
			for len(ch) > 0 {
				<-ch
			}
			assert.Nil(c.LeadershipTransfer())
			// receive the leadership lost signal
			assert.False(<-ch)
			assert.False(c.IsLeader())
		}
	}

	// make sure the leader is chanded and the same on all the peers
	var leader string
	timeout = time.After(10 * time.Second)
loop1:
	for {
		select {
		case <-timeout:
			t.Logf("Election timeout")
			t.Fail()
		default:
			for _, c := range clusters {
				l := c.LeaderAddress()
				if leader == "" || leader != l {
					leader = l
					continue loop1
				}
			}
			break loop1
		}
	}

	n := 0
	for _, c := range clusters {
		assert.NotEqual(addr1, c.LeaderAddress())
		if c.IsLeader() {
			n++
		}
	}
	assert.Equal(1, n)

	// Shutdown
	for _, c := range clusters {
		assert.Nil(c.Shutdown())
	}
}
