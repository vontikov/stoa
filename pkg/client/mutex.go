package client

import (
	"context"
	"sync"

	"google.golang.org/grpc/metadata"

	"github.com/vontikov/stoa/pkg/pb"
)

type mutex struct {
	base
	clientID []byte

	mu    sync.RWMutex // protects following fields
	watch chan *pb.MutexStatus
}

func (c *client) newMutex(name string) *mutex {
	return &mutex{
		base:     c.createBase(name),
		clientID: c.id,
	}
}

// TTryLock implements Mutex.TryLock.
func (m *mutex) TryLock(ctx context.Context, opts ...CallOption) (r bool, err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	msg := pb.ClientId{
		EntityName: m.name,
		Id:   m.clientID,
	}

	o, err := m.handle.MutexTryLock(ctx, &msg, m.callOptions...)
	if err == nil {
		r = o.Ok
		return
	}
	if m.failFast {
		return
	}

	err = retry(ctx,
		func() error {
			o, err := m.handle.MutexTryLock(ctx, &msg, m.callOptions...)
			if err == nil {
				r = o.Ok
				return nil
			}
			return err
		},
		m.retryTimeout, m.idleStrategy,
	)
	return
}

// Unlock implements Mutex.Unlock.
func (m *mutex) Unlock(ctx context.Context, opts ...CallOption) (r bool, err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	msg := pb.ClientId{
		EntityName: m.name,
		Id:   m.clientID,
	}

	o, err := m.handle.MutexUnlock(ctx, &msg, m.callOptions...)
	if err == nil {
		r = o.Ok
		return
	}
	if m.failFast {
		return
	}

	err = retry(ctx,
		func() error {
			o, err := m.handle.MutexUnlock(ctx, &msg, m.callOptions...)
			if err == nil {
				r = o.Ok
				return nil
			}
			return err
		},
		m.retryTimeout, m.idleStrategy,
	)
	return
}

// Watch implements Mutex.Watch.
func (m *mutex) Watch() <-chan *pb.MutexStatus {
	m.mu.Lock()
	if m.watch == nil {
		m.watch = make(chan *pb.MutexStatus, watchChanSize)
	}
	m.mu.Unlock()
	return m.watch
}
