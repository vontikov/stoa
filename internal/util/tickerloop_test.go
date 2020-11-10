package util

import (
	"testing"

	"context"
	"sync/atomic"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTickerLoopCallOnce(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var c int32
	f := func() error {
		time.Sleep(200 * time.Millisecond)
		atomic.AddInt32(&c, 1)
		return nil
	}

	err := TickerLoop(ctx, 10*time.Millisecond, f, nil)
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Equal(t, int32(1), c)
}

func TestEventLoopCallNotOnce(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var c int32
	f := func() error {
		atomic.AddInt32(&c, 1)
		return nil
	}

	err := TickerLoop(ctx, 10*time.Millisecond, f, nil)
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Less(t, int32(1), c)
}
