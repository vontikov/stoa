package util

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type ErrorHandler func(error) bool

func TickerLoop(ctx context.Context, p time.Duration, f func() error, h ErrorHandler) error {
	ticker := time.NewTicker(p)
	var flag int32
	var wg sync.WaitGroup
	defer wg.Wait()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if !atomic.CompareAndSwapInt32(&flag, 0, 1) {
				continue
			}

			wg.Add(1)
			go func() {
				err := f()
				if h != nil && !h(err) {
					return
				}
				atomic.StoreInt32(&flag, 0)
				wg.Done()
			}()
		}
	}
}
