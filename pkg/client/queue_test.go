package client

import (
	"testing"

	"bytes"
	"sync"

	"context"
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/vontikov/stoa/internal/test"
)

func TestQueue(t *testing.T) {
	const (
		clusterSize = 3
		basePort    = 2000
		queueName   = "queue"
		max         = 100
	)

	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, bootstrap, err := test.StartCluster(ctx, basePort, clusterSize)
	assert.Nil(err)

	client, err := New(ctx, WithBootstrap(bootstrap))
	assert.Nil(err)

	queue := client.Queue(queueName)
	sz, err := queue.Size(ctx)
	assert.Nil(err)
	assert.Equal(uint32(0), sz, "Queue should be empty")

	for i := 0; i < max; i++ {
		var wg sync.WaitGroup
		wg.Add(1)
		go func(i int) {
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "Message #%d", i)
			err = queue.Offer(ctx, buf.Bytes())
			assert.Nil(err, "Offer() should return nil")

			wg.Done()
		}(i)
		wg.Wait()
	}

	sz, err = queue.Size(ctx)
	assert.Nil(err)
	assert.Equal(uint32(max), sz, "All messages should be added")

	for i := 0; i < max/2; i++ {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Message #%d", i)

		v, err := queue.Peek(ctx)
		assert.Nil(err, "Peek() should return nil")
		assert.Equal(buf.Bytes(), v)

		sz, err = queue.Size(ctx)
		assert.Nil(err)
		assert.Equal(uint32(max-i), sz, "Peek() should not remove message")

		v, err = queue.Poll(ctx)
		assert.Nil(err, "Poll() should return nil")
		assert.Equal(buf.Bytes(), v)

		sz, err = queue.Size(ctx)
		assert.Nil(err)
		assert.Equal(uint32(max-i-1), sz, "Poll() should remove message")
	}

	err = queue.Clear(ctx)
	assert.Nil(err, "Clear() should return nil")

	sz, err = queue.Size(ctx)
	assert.Nil(err)
	assert.Equal(uint32(0), sz, "Queue should be empty")
}
