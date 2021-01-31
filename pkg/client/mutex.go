package client

import (
	"context"

	"github.com/vontikov/stoa/pkg/pb"
	"google.golang.org/grpc/metadata"
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

func (m *mutex) TryLock(ctx context.Context, opts ...CallOption) (r bool, err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	msg := pb.Id{
		Name: m.name,
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

func (m *mutex) Unlock(ctx context.Context, opts ...CallOption) (r bool, err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	msg := pb.Id{
		Name: m.name,
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
