package client

import (
	"context"
	"encoding/binary"

	"google.golang.org/grpc/metadata"

	"github.com/vontikov/stoa/pkg/pb"
)

type queue struct {
	base
}

func (c *client) newQueue(name string) *queue {
	return &queue{base: c.createBase(name)}
}

// Size implements Queue.Size.
func (q *queue) Size(ctx context.Context, opts ...CallOption) (r uint32, err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	m := pb.Entity{EntityName: q.name}
	v, err := q.handle.QueueSize(ctx, &m, q.callOptions...)
	if err == nil {
		r = binary.LittleEndian.Uint32(v.Value)
		return
	}
	if q.failFast {
		return
	}

	err = retry(ctx,
		func() error {
			v, err := q.handle.QueueSize(ctx, &m, q.callOptions...)
			if err != nil {
				return err
			}
			r = binary.LittleEndian.Uint32(v.Value)
			return nil
		},
		q.retryTimeout, q.idleStrategy,
	)

	return
}

// Clear implements Queue.Clear.
func (q *queue) Clear(ctx context.Context, opts ...CallOption) (err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	m := pb.Entity{EntityName: q.name}
	_, err = q.handle.QueueClear(ctx, &m, q.callOptions...)
	if err == nil || q.failFast {
		return
	}

	err = retry(ctx,
		func() error {
			_, err := q.handle.QueueClear(ctx, &m, q.callOptions...)
			return err
		},
		q.retryTimeout, q.idleStrategy,
	)
	return
}

// Offer implements Queue.Offer.
func (q *queue) Offer(ctx context.Context, e []byte, opts ...CallOption) (err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	m := pb.Value{
		EntityName: q.name,
		Value:      e,
	}
	_, err = q.handle.QueueOffer(ctx, &m, q.callOptions...)
	if err == nil || q.failFast {
		return
	}

	err = retry(ctx,
		func() error {
			_, err := q.handle.QueueOffer(ctx, &m, q.callOptions...)
			return err
		},
		q.retryTimeout, q.idleStrategy,
	)
	return
}

// Poll implements Queue.Poll.
func (q *queue) Poll(ctx context.Context, opts ...CallOption) (r []byte, err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	m := pb.Entity{EntityName: q.name}
	v, err := q.handle.QueuePoll(ctx, &m, q.callOptions...)
	if err == nil {
		r = v.Value
		return
	}
	if q.failFast {
		return nil, err
	}

	err = retry(ctx,
		func() error {
			v, err := q.handle.QueuePoll(ctx, &m, q.callOptions...)
			if err != nil {
				return err
			}
			r = v.Value
			return nil
		},
		q.retryTimeout, q.idleStrategy,
	)
	return
}

// Peek implements Queue.Peek.
func (q *queue) Peek(ctx context.Context, opts ...CallOption) (r []byte, err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	m := pb.Entity{EntityName: q.name}
	v, err := q.handle.QueuePeek(ctx, &m, q.callOptions...)
	if err == nil {
		r = v.Value
		return
	}
	if q.failFast {
		return nil, err
	}

	err = retry(ctx,
		func() error {
			v, err := q.handle.QueuePeek(ctx, &m, q.callOptions...)
			if err != nil {
				return err
			}
			r = v.Value
			return nil
		},
		q.retryTimeout, q.idleStrategy,
	)
	return
}

