package client

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/grpc/metadata"

	concurrent "github.com/vontikov/go-concurrent"
	"github.com/vontikov/stoa/internal/common"
)

type retryFunc func() error

// TimeNowInMillis returns current time in milliseconds.
var TimeNowInMillis = func() int64 { return time.Now().UnixNano() / int64(time.Millisecond) }

func retry(ctx context.Context, f retryFunc, t time.Duration, s concurrent.IdleStrategy) error {
	ctx, cancel := context.WithTimeout(ctx, t)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := f()
			if err == nil {
				return nil
			}
			s.Idle()
		}
	}
}

// MetadataFromCallOptions creates gRPC metadata from call options.
func MetadataFromCallOptions(opts ...CallOption) metadata.MD {
	options := &callOptions{}
	for _, o := range opts {
		o(options)
	}
	md := metadata.MD{}
	if options.barrierEnabled {
		md[common.MetaKeyBarrier] = []string{strconv.Itoa(int(options.barrier / time.Millisecond))}
	}
	if options.ttlEnabled {
		ttl := TimeNowInMillis() + int64(options.ttl/time.Millisecond)
		md[common.MetaKeyTTL] = []string{strconv.Itoa(int(ttl))}
	}
	return md
}
