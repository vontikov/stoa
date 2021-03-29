package client

import (
	"context"
	"crypto/rand"
	"math/big"
	"strconv"
	"time"

	"google.golang.org/grpc/metadata"

	concurrent "github.com/vontikov/go-concurrent"
	"github.com/vontikov/stoa/pkg/common"
)

// TimeNowInMillis returns current time in milliseconds.
var TimeNowInMillis = func() int64 { return time.Now().UnixNano() / int64(time.Millisecond) }

type retryFunc func() error
type idleFunc func()

func retry(ctx context.Context, fn retryFunc, timeout time.Duration, idle concurrent.IdleStrategy, onIdleFn ...idleFunc) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := fn()
			if err == nil {
				return nil
			}
			idle.Idle()
			for _, f := range onIdleFn {
				f()
			}
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

var idSrc = []byte("abcdef0123456789")

// genID returns random ID of length n.
func genID(n int) []byte {
	max := big.NewInt(int64(len(idSrc)))
	b := make([]byte, n)
	for i := range b {
		idx, _ := rand.Int(rand.Reader, max)
		b[i] = idSrc[idx.Uint64()]
	}
	return b
}
