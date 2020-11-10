package client

import (
	"context"

	cc "github.com/vontikov/go-concurrent"

	"github.com/vontikov/stoa/pkg/pb"
	"google.golang.org/grpc/metadata"
)

type mutex struct {
	base
	clientID string
	watchers cc.List
}

func newMutex(name string, cfg *options, handle pb.StoaClient) *mutex {
	return &mutex{
		base:     createBase(name, cfg, handle),
		clientID: cfg.id,
		watchers: cc.NewSynchronizedList(0),
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

func (m *mutex) Watch(w MutexWatch) {
	m.watchers.Add(w)
}

func (m *mutex) Unwatch(w MutexWatch) {
	m.watchers.Remove(w, func(l, r interface{}) bool { return l == r })
}
