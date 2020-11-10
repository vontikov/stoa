package client

import (
	"testing"

	"context"
	"errors"
	"time"

	"github.com/stretchr/testify/assert"
	concurrent "github.com/vontikov/go-concurrent"
	"github.com/vontikov/stoa/internal/common"
)

func TestRetryFunc(t *testing.T) {
	const timeout = 100 * time.Millisecond
	const idleTimeout = timeout / 10

	is := concurrent.NewSleepingIdleStrategy(idleTimeout)

	t.Run("Should execute",
		func(t *testing.T) {
			err := retry(
				context.Background(),
				func() error { return nil },
				timeout,
				is)
			assert.Nil(t, err)
		})

	t.Run("Should backoff",
		func(t *testing.T) {
			n := 0
			err := retry(
				context.Background(),
				func() error {
					if n == 0 {
						n++
						return errors.New("an error")
					}
					return nil
				},
				timeout,
				is)
			assert.Equal(t, 1, n)
			assert.Nil(t, err)
		})

	t.Run("Should reach retry timeout",
		func(t *testing.T) {
			err := retry(
				context.Background(),
				func() error { return errors.New("an error") },
				timeout,
				is)
			assert.Equal(t, context.DeadlineExceeded, err)
		})
}

func TestMetadataFromCallOptions(t *testing.T) {
	// Make time-based calls testable
	TimeNowInMillis = func() int64 { return 0 }

	t.Run("Should enable barrier",
		func(t *testing.T) {
			md := MetadataFromCallOptions(WithBarrier())
			v := md[common.MetaKeyBarrier][0]
			assert.Equal(t, "0", v)
		})

	t.Run("Should enable waiting barrier",
		func(t *testing.T) {
			md := MetadataFromCallOptions(WithWaitingBarrier(5 * time.Second))
			v := md[common.MetaKeyBarrier][0]
			assert.Equal(t, "5000", v, "Should be 5000 milliseconds")
		})

	t.Run("Should enable message TTL",
		func(t *testing.T) {
			md := MetadataFromCallOptions(WithTTL(42 * time.Second))
			v := md[common.MetaKeyTTL]
			assert.Equal(t, "42000", v[0], "Should be 42000 milliseconds")
		})
}
