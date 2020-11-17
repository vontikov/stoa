package client

/*
import (
	"testing"

	"bytes"
	"sync"

	"context"
	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	const (
		queueName = "queue"
		max       = 100
	)

	defer runTestNode(t)()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := New(
		WithContext(ctx),
		WithDiscoveryIP(testDiscoveryIP),
		WithDiscoveryPort(testDiscoveryPort),
	)
	assert.Nil(t, err)
	assert.NotNil(t, client)

	queue := client.Queue(queueName)
	sz, err := queue.Size(ctx)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), sz, "Queue should be empty")

	for i := 0; i < max; i++ {
		var wg sync.WaitGroup
		wg.Add(1)
		go func(i int) {
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "Message #%d", i)
			err = queue.Offer(ctx, buf.Bytes())
			assert.Nil(t, err, "Offer() should return nil")

			wg.Done()
		}(i)
		wg.Wait()
	}

	sz, err = queue.Size(ctx)
	assert.Nil(t, err)
	assert.Equal(t, uint32(max), sz, "All messages should be added")

	for i := 0; i < max/2; i++ {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Message #%d", i)

		v, err := queue.Peek(ctx)
		assert.Nil(t, err, "Peek() should return nil")
		assert.Equal(t, buf.Bytes(), v)

		sz, err = queue.Size(ctx)
		assert.Nil(t, err)
		assert.Equal(t, uint32(max-i), sz, "Peek() should not remove message")

		v, err = queue.Poll(ctx)
		assert.Nil(t, err, "Poll() should return nil")
		assert.Equal(t, buf.Bytes(), v)

		sz, err = queue.Size(ctx)
		assert.Nil(t, err)
		assert.Equal(t, uint32(max-i-1), sz, "Poll() should remove message")
	}

	err = queue.Clear(ctx)
	assert.Nil(t, err, "Clear() should return nil")

	sz, err = queue.Size(ctx)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), sz, "Queue should be empty")
}
*/
