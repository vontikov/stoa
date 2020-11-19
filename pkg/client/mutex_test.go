package client

import (
	"sync"
	"testing"

	"context"
	"fmt"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/test"
)

func TestMutex(t *testing.T) {
	const (
		muxName  = "mx"
		basePort = 2900
	)

	assert := assert.New(t)

	addr1 := fmt.Sprintf("127.0.0.1:%d", basePort+1)
	addr2 := fmt.Sprintf("127.0.0.1:%d", basePort+2)
	addr3 := fmt.Sprintf("127.0.0.1:%d", basePort+3)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	test.RunTestCluster(ctx, t, basePort)

	peers := strings.Join([]string{addr1, addr2, addr3}, cluster.PeerListSep)
	client0, err := New(WithContext(ctx), WithPeers(peers))
	assert.Nil(err)

	mx0 := client0.Mutex(muxName)
	ok, err := mx0.TryLock(ctx)
	assert.Nil(err)
	assert.True(ok)

	ok, err = mx0.TryLock(ctx)
	assert.Nil(err)
	assert.False(ok)

	client1, err := New(WithContext(ctx), WithPeers(peers))
	assert.Nil(err)

	mx1 := client1.Mutex(muxName)
	ok, err = mx1.TryLock(ctx)
	assert.Nil(err)
	assert.False(ok)

	ok, err = mx0.Unlock(ctx)
	assert.Nil(err)
	assert.True(ok)
	ok, err = mx0.Unlock(ctx)
	assert.Nil(err)
	assert.False(ok)

	ok, err = mx1.Unlock(ctx)
	assert.Nil(err)
	assert.False(ok)

	cancel()
}

func TestMutexWatch(t *testing.T) {
	const (
		max      = 10
		basePort = 3000
	)

	assert := assert.New(t)

	addr1 := fmt.Sprintf("127.0.0.1:%d", basePort+1)
	addr2 := fmt.Sprintf("127.0.0.1:%d", basePort+2)
	addr3 := fmt.Sprintf("127.0.0.1:%d", basePort+3)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	test.RunTestCluster(ctx, t, basePort)

	mxName := func(i int) string { return fmt.Sprintf("mx-%d", i) }

	clients := make(chan Client, max)
	peers := strings.Join([]string{addr1, addr2, addr3}, cluster.PeerListSep)

	var wg sync.WaitGroup
	for i := 0; i < max; i++ {
		wg.Add(1)
		id := fmt.Sprintf("id-%d", i)
		go func() {
			client, err := New(WithID(id), WithContext(ctx), WithPeers(peers))
			assert.Nil(err)

			for j := 0; j < max; j++ {
				name := mxName(j)
				mx := client.Mutex(name)
				w := struct{ MutexWatchProto }{}
				w.Callback = func(n string, v bool) {
					assert.Equal(name, n)
				}
				mx.Watch(w)
			}

			clients <- client
			wg.Done()
		}()
	}
	wg.Wait()
	close(clients)

	var i int
	for c := range clients {
		c.Mutex(mxName(i)).TryLock(context.Background())
		i++
	}
}
