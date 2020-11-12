package cluster

import (
	"testing"

	"context"
	"sync"

	"github.com/stretchr/testify/assert"
)

func TestFsmQueue(t *testing.T) {
	const queueName = "test-queue"
	const max = 100

	f := NewFSM(context.Background())

	var wg sync.WaitGroup
	for i := 0; i < max; i++ {
		wg.Add(1)
		go func(n int) {
			q := queue(f, queueName)
			assert.NotNil(t, q)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
