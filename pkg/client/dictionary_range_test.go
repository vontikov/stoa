// +build all

package client

import (
	"testing"

	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/gateway"
	"github.com/vontikov/stoa/internal/test"
)

func TestDictionaryRange(t *testing.T) {
	const (
		basePort = 2900
		dictName = "dict"
		max      = 10000
	)

	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clusters := test.RunTestCluster(ctx, t, basePort)

	addr1 := fmt.Sprintf("127.0.0.1:%d", basePort+1)
	addr2 := fmt.Sprintf("127.0.0.1:%d", basePort+2)
	addr3 := fmt.Sprintf("127.0.0.1:%d", basePort+3)
	peers := strings.Join([]string{addr1, addr2, addr3}, cluster.PeerListSep)

	client, err := New(WithContext(ctx), WithPeers(peers))
	assert.Nil(err)

	dict := client.Dictionary(dictName)

	// populate
	for i := 0; i < max; i++ {
		k := strconv.Itoa(i)
		v := strconv.Itoa(i + max)
		ok, err := dict.PutIfAbsent(ctx, []byte(k), []byte(v))
		assert.Nil(err)
		assert.True(ok)
	}

	tctx, tcancel := context.WithTimeout(ctx, 30*time.Second)
	defer tcancel()

	// transfer leadership
	var trs int32
	go func() {
		t := time.Tick(5 * time.Second)
		for {
			select {
			case <-tctx.Done():
				return
			case <-t:
				for _, c := range clusters {
					if c.IsLeader() {
						err := c.LeadershipTransfer()
						assert.Nil(err)
						atomic.AddInt32(&trs, 1)
					}
				}
			}
		}
	}()

	var pos, neg int32

	testFunc := func() {
		t.Logf("start")
		m := make(map[int]int)
		kvChan, errChan := dict.Range(ctx)

	loop:
		for {
			select {
			case err := <-errChan:
				// happens while the new leader is being elected
				assert.True(errors.Is(err, gateway.ErrNotLeader))
				t.Logf("warn: %s", err)
				atomic.AddInt32(&neg, 1)
				return
			case kv := <-kvChan:
				if kv == nil {
					break loop
				}
				k, err := strconv.Atoi(string(kv[0]))
				assert.Nil(err)
				v, err := strconv.Atoi(string(kv[1]))
				assert.Nil(err)
				m[k] = v
			}
		}

		assert.Equal(max, len(m))
		for i := 0; i < max; i++ {
			assert.Equal(i+max, m[i])
		}

		atomic.AddInt32(&pos, 1)
		t.Logf("done")
	}

loop:
	for {
		select {
		case <-tctx.Done():
			break loop
		default:
			var wg sync.WaitGroup
			for i := 0; i < 3; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					testFunc()
					time.Sleep(1 * time.Second)
				}()
			}
			wg.Wait()
		}
	}

	t.Logf("trs=%d", trs)
	assert.Less(int32(3), trs)

	t.Logf("pos=%d", pos)
	assert.Less(int32(60), pos)

	t.Logf("neg=%d", neg)
	assert.Greater(int32(20), neg)
}
