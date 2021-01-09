package client

import (
	"context"
	"encoding/binary"
	"io"

	"google.golang.org/grpc/metadata"

	"github.com/vontikov/stoa/pkg/pb"
)

type dictionary struct {
	base
}

func (c *client) newDictionary(name string) *dictionary {
	return &dictionary{base: c.createBase(name)}
}

// Size implements Dictionary.Size.
func (d *dictionary) Size(ctx context.Context, opts ...CallOption) (r uint32, err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	m := pb.Name{Name: d.name}
	v, err := d.handle.DictionarySize(ctx, &m, d.callOptions...)
	if err == nil {
		r = binary.LittleEndian.Uint32(v.Value)
		return
	}
	if d.failFast {
		return
	}

	err = retry(ctx,
		func() error {
			v, err := d.handle.DictionarySize(ctx, &m, d.callOptions...)
			if err != nil {
				return err
			}
			r = binary.LittleEndian.Uint32(v.Value)
			return nil
		},
		d.retryTimeout, d.idleStrategy,
	)
	return
}

// Clear implements Dictionary.Clear.
func (d *dictionary) Clear(ctx context.Context, opts ...CallOption) (err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	m := pb.Name{Name: d.name}
	_, err = d.handle.DictionaryClear(ctx, &m, d.callOptions...)
	if err == nil || d.failFast {
		return
	}

	err = retry(ctx,
		func() error {
			_, err := d.handle.DictionaryClear(ctx, &m, d.callOptions...)
			return err
		},
		d.retryTimeout, d.idleStrategy,
	)
	return
}

// Put implements Dictionary.Put.
func (d *dictionary) Put(ctx context.Context, k, v []byte, opts ...CallOption) (r []byte, err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	m := pb.KeyValue{
		Name: d.name,
		Key:  k, Value: v,
	}
	o, err := d.handle.DictionaryPut(ctx, &m, d.callOptions...)
	if err == nil {
		r = o.Value
		return
	}
	if d.failFast {
		return
	}

	err = retry(ctx,
		func() error {
			o, err := d.handle.DictionaryPut(ctx, &m, d.callOptions...)
			if err == nil {
				r = o.Value
				return nil
			}
			return err
		},
		d.retryTimeout, d.idleStrategy,
	)
	return
}

// PutIfAbsent implements Dictionary.PutIfAbsent.
func (d *dictionary) PutIfAbsent(ctx context.Context, k, v []byte, opts ...CallOption) (r bool, err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	m := pb.KeyValue{
		Name: d.name,
		Key:  k, Value: v,
	}
	o, err := d.handle.DictionaryPutIfAbsent(ctx, &m, d.callOptions...)
	if err == nil {
		r = o.Ok
		return
	}
	if d.failFast {
		return
	}

	err = retry(ctx,
		func() error {
			o, err := d.handle.DictionaryPutIfAbsent(ctx, &m, d.callOptions...)
			if err == nil {
				r = o.Ok
				return nil
			}
			return err
		},
		d.retryTimeout, d.idleStrategy,
	)
	return
}

// Get implements Dictionary.Get.
func (d *dictionary) Get(ctx context.Context, k []byte, opts ...CallOption) (r []byte, err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	m := pb.Key{
		Name: d.name,
		Key:  k,
	}
	o, err := d.handle.DictionaryGet(ctx, &m, d.callOptions...)
	if err == nil {
		r = o.Value
		return
	}
	if d.failFast {
		return
	}

	err = retry(ctx,
		func() error {
			o, err := d.handle.DictionaryGet(ctx, &m, d.callOptions...)
			if err == nil {
				r = o.Value
				return nil
			}
			return err
		},
		d.retryTimeout, d.idleStrategy,
	)
	return
}

// Remove implements Dictionary.Remove.
func (d *dictionary) Remove(ctx context.Context, k []byte, opts ...CallOption) (err error) {
	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	m := pb.Key{
		Name: d.name,
		Key:  k,
	}
	_, err = d.handle.DictionaryRemove(ctx, &m, d.callOptions...)
	if err == nil {
		return
	}
	if d.failFast {
		return
	}

	err = retry(ctx,
		func() error {
			_, err := d.handle.DictionaryRemove(ctx, &m, d.callOptions...)
			return err
		},
		d.retryTimeout, d.idleStrategy,
	)
	return
}

func (d *dictionary) Range(ctx context.Context, opts ...CallOption) (<-chan [][]byte, <-chan error) {
	rc := make(chan [][]byte)
	ec := make(chan error, 1)

	if len(opts) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, MetadataFromCallOptions(opts...))
	}

	go func() {
		m := pb.Name{Name: d.name}
		stream, err := d.handle.DictionaryRange(ctx, &m, d.callOptions...)
		if err != nil {
			if d.failFast {
				ec <- err
				return
			}
			if err := retry(ctx, func() error {
				var err error
				stream, err = d.handle.DictionaryRange(ctx, &m, d.callOptions...)
				return err
			},
				d.retryTimeout, d.idleStrategy); err != nil {
				ec <- err
				return
			}
		}
		for {
			kv, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					close(rc)
					return
				}
				ec <- err
				return
			}
			rc <- [][]byte{kv.Key, kv.Value}
		}
	}()

	return rc, ec
}
