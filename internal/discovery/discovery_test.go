package discovery

import (
	"context"
	"net"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/stretchr/testify/assert"
)

func TestDiscovery(t *testing.T) {
	assert := assert.New(t)

	const (
		ip   = "224.0.0.1"
		port = 1235
		msg  = "Hello, World!"
		max  = 10
	)

	Duration(50 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	n := 0
	handler := func(src *net.UDPAddr, m []byte) error {
		assert.Equal(msg, string(m))
		n++
		if n == max {
			cancel()
		}
		return nil
	}

	s, err := NewSender(ip, port, func() ([]byte, error) { return []byte(msg), nil })
	assert.Nil(err)

	r, err := NewReceiver(ip, port, handler)
	assert.Nil(err)

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error { return s.Run(gctx) })
	g.Go(func() error { return r.Run(gctx) })
	err = g.Wait()
	assert.Equal(context.Canceled, err)
}
