package client

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vontikov/stoa/internal/test"
	"github.com/vontikov/stoa/pkg/logging"
	"golang.org/x/sync/errgroup"
)

func TestClientStart(t *testing.T) {
	const (
		clusterSize = 3
		basePort    = 2100
		dictName    = "dict"
	)

	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, bootstrap, err := test.StartCluster(ctx, basePort, clusterSize)
	assert.Nil(err)

	client, err := New(ctx, WithBootstrap(bootstrap))
	assert.Nil(err)

	dict := client.Dictionary(dictName)
	sz, err := dict.Size(ctx)
	assert.Nil(err)
	assert.Equal(uint32(0), sz)
}

func TestMutexExpiration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	logging.SetLevel("trace")

	const (
		clusterSize = 3
		basePort    = 2200
		mutexName   = "test-mutex"
	)

	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, bootstrap, err := test.StartCluster(ctx, basePort, clusterSize)
	assert.Nil(err)

	g, ctx := errgroup.WithContext(ctx)
	runChan := make(chan bool)
	payload := []byte("some payload")

	g.Go(func() error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		client, err := New(ctx, WithBootstrap(bootstrap))
		if err != nil {
			return err
		}

		mx := client.Mutex(mutexName)

		<-runChan
		r, p, err := mx.TryLock(ctx, nil)
		if err != nil {
			return err
		}
		assert.False(r, "should be locked by another client")
		assert.Equal(payload, p, "the payload should be provided by the holder")
		return nil
	})

	// make sure the first client is fully initialized
	time.Sleep(200 * time.Millisecond)

	g.Go(func() error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		client, err := New(ctx, WithBootstrap(bootstrap))
		if err != nil {
			return err
		}
		mx := client.Mutex(mutexName)

		r, p, err := mx.TryLock(ctx, payload)
		if err != nil {
			return err
		}
		assert.True(r)
		assert.Equal([]byte(nil), p, "the payload should be empty")
		runChan <- true
		return nil
	})

	err = g.Wait()
	assert.Nil(err)
}
