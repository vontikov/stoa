package client

/*
import (
	"testing"

	"fmt"
	"sync"

	"context"

	"github.com/stretchr/testify/assert"
)

func TestMutex(t *testing.T) {
	const muxName = "mx"

	defer runTestNode(t)()

	ctx, cancel := context.WithCancel(context.Background())

	client0, err := New(
		WithContext(ctx),
		WithDiscoveryIP(testDiscoveryIP),
		WithDiscoveryPort(testDiscoveryPort),
	)
	assert.Nil(t, err)

	mx0 := client0.Mutex(muxName)
	ok, err := mx0.TryLock(ctx)
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = mx0.TryLock(ctx)
	assert.Nil(t, err)
	assert.False(t, ok)

	client1, err := New(
		WithContext(ctx),
		WithDiscoveryIP(testDiscoveryIP),
		WithDiscoveryPort(testDiscoveryPort),
	)
	assert.Nil(t, err)

	mx1 := client1.Mutex(muxName)
	ok, err = mx1.TryLock(ctx)
	assert.Nil(t, err)
	assert.False(t, ok)

	ok, err = mx0.Unlock(ctx)
	assert.Nil(t, err)
	assert.True(t, ok)
	ok, err = mx0.Unlock(ctx)
	assert.Nil(t, err)
	assert.False(t, ok)

	ok, err = mx1.Unlock(ctx)
	assert.Nil(t, err)
	assert.False(t, ok)

	cancel()
}

func TestMutexWatch(t *testing.T) {
	const max = 10

	defer runTestNode(t)()

	mxName := func(i int) string { return fmt.Sprintf("mx-%d", i) }

	ctx, cancel := context.WithCancel(context.Background())
	clients := make(chan Client, max)

	var wg sync.WaitGroup
	for i := 0; i < max; i++ {
		wg.Add(1)
		id := fmt.Sprintf("id-%d", i)
		go func() {
			client, err := New(
				WithContext(ctx),
				WithDiscoveryIP(testDiscoveryIP),
				WithDiscoveryPort(testDiscoveryPort),
				WithLogLevel(testLogLevel),
				WithID(id),
			)
			assert.Nil(t, err)

			for j := 0; j < max; j++ {
				name := mxName(j)
				mx := client.Mutex(name)
				w := struct{ MutexWatchProto }{}
				w.Callback = func(n string, v bool) {
					assert.Equal(t, name, n)
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

	cancel()
}
*/
