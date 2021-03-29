package client

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/vontikov/stoa/pkg/pb"
)

type mutex struct {
	base
	clientID []byte
}

func (c *client) newMutex(name string) *mutex {
	return &mutex{
		base:     c.createBase(name),
		clientID: c.id,
	}
}

// TTryLock implements Mutex.TryLock.
func (m *mutex) TryLock(ctx context.Context, payload []byte, opts ...CallOption) (r bool, p []byte, err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	msg := pb.ClientId{
		EntityName: m.name,
		Id:         m.clientID,
		Payload:    payload,
	}

	o, err := m.handle.MutexTryLock(ctx, &msg, m.callOptions...)
	if err == nil {
		r = o.Ok
		p = o.Payload
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
				p = o.Payload
				return nil
			}
			return err
		},
		m.retryTimeout, m.idleStrategy,
	)
	return
}

// Unlock implements Mutex.Unlock.
func (m *mutex) Unlock(ctx context.Context, opts ...CallOption) (r bool, payload []byte, err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	msg := pb.ClientId{
		EntityName: m.name,
		Id:         m.clientID,
	}

	o, err := m.handle.MutexUnlock(ctx, &msg, m.callOptions...)
	if err == nil {
		r = o.Ok
		payload = o.Payload
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
				payload = o.Payload
				return nil
			}
			return err
		},
		m.retryTimeout, m.idleStrategy,
	)
	return
}
