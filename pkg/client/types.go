package client

import (
	"context"
	"time"

	cc "github.com/vontikov/go-concurrent"

	"github.com/vontikov/stoa/pkg/pb"
	"google.golang.org/grpc"
)

// Client defines Stoa client methods.
type Client interface {
	Wait()
	Queue(n string) Queue
	Dictionary(n string) Dictionary
	Mutex(n string) Mutex
}

// CallOption adds an option to the method call.
type CallOption func(*callOptions)

// WithTTL is used to set the element's TTL.
func WithTTL(d time.Duration) CallOption {
	return func(o *callOptions) {
		o.ttlEnabled = true
		o.ttl = d
	}
}

// WithWaitingBarrier is used to block a call until all preceding have been
// applied to the Raft Cluster.
// An optional timeout limits the amount of time we wait for the command to be started.
func WithWaitingBarrier(d time.Duration) CallOption {
	return func(o *callOptions) {
		o.barrierEnabled = true
		o.barrier = d
	}
}

// WithBarrier is used to block a call until all preceding have been
// applied to the Raft Cluster.
func WithBarrier() CallOption {
	return func(o *callOptions) {
		o.barrierEnabled = true
	}
}

// Collection defines generic collection
type Collection interface {

	// Size returns the collection size
	Size(ctx context.Context, opts ...CallOption) (uint32, error)

	// Clear clears the collection
	Clear(ctx context.Context, opts ...CallOption) error
}

// Queue defines ordered in FIFO manner collection of elements which may contain
// duplicates
type Queue interface {
	Collection

	// Offer inserts the element e into the queue
	Offer(ctx context.Context, e []byte, opts ...CallOption) error

	// Poll retrieves and removes the head of the queue; returns nil if the queue is empty
	Poll(ctx context.Context, opts ...CallOption) ([]byte, error)

	// Peek retrieves, but does not remove, the head of the queue; returns nil if the queue is empty
	Peek(ctx context.Context, opts ...CallOption) ([]byte, error)
}

// Dictionary represents a collection of key-value pairs.
type Dictionary interface {
	Collection

	// Put puts a new key-value pair into the dictionary.
	// If the key already exists overwrites the existing value with the new one.
	// Returns the previous value associated with key, or nil if there was no mapping
	// for key.
	Put(ctx context.Context, k, v []byte, opts ...CallOption) ([]byte, error)

	// PutIfAbsent puts the key-value pair (and returns true)
	// only if the key is absent, otherwise it returns false.
	PutIfAbsent(ctx context.Context, k, v []byte, opts ...CallOption) (bool, error)

	// Get returns the value specified by the key if the key-value pair is
	// present, othervise returns nil.
	Get(ctx context.Context, k []byte, opts ...CallOption) ([]byte, error)

	// Remove removes the key-value pair specified by the key k from the map
	// if it is present.
	Remove(ctx context.Context, k []byte, opts ...CallOption) error

	// Range returns the channel kv that will send all key-value pairs containing in
	// the dictionary.
	Range(ctx context.Context, opts ...CallOption) (kv <-chan [][]byte, err <-chan error)
}

type MutexWatcher interface {
	Apply(string, bool)
}

type MutexWatchProto struct {
	Callback func(string, bool)
}

func (m MutexWatchProto) Apply(n string, v bool) { m.Callback(n, v) }

type Mutex interface {
	TryLock(ctx context.Context, opts ...CallOption) (bool, error)
	Unlock(ctx context.Context, opts ...CallOption) (bool, error)
	Watch(MutexWatcher)
	Unwatch(MutexWatcher)
}

type callOptions struct {
	barrierEnabled bool
	barrier        time.Duration
	ttlEnabled     bool
	ttl            time.Duration
}

type base struct {
	name         string
	callOptions  []grpc.CallOption
	failFast     bool
	handle       pb.StoaClient
	idleStrategy cc.IdleStrategy
	retryTimeout time.Duration
}

func createBase(name string, cfg *options, handle pb.StoaClient) base {
	return base{
		name:         name,
		handle:       handle,
		callOptions:  cfg.callOptions,
		failFast:     cfg.failFast,
		idleStrategy: cfg.idleStrategy,
		retryTimeout: cfg.retryTimeout,
	}
}
