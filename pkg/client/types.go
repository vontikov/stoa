package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	cc "github.com/vontikov/go-concurrent"

	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
	"google.golang.org/grpc"
)

// ErrDeadlineExceeded is the error returned if a timeout is specified.
var ErrDeadlineExceeded = errors.New("deadline exceeded")

// Collection is a generic collection
type Collection interface {

	// Name returns the Collection name
	Name() string

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

	// Watch returns the Queue notification channel.
	Watch() <-chan interface{}
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

// Mutex is a mutual exclusion lock.
type Mutex interface {
	TryLock(ctx context.Context, opts ...CallOption) (bool, error)
	Unlock(ctx context.Context, opts ...CallOption) (bool, error)
}

const (
	defaultDialTimeout     = 30000 * time.Millisecond
	defaultKeepAlivePeriod = 1000 * time.Millisecond
	defaultRetryTimeout    = 15000 * time.Millisecond
	defaultClientIDSize    = 16
	defaultLoggerName      = "stoa-client"
	peerListSep            = ","
)

// Option is a function applied to an options to change the options' default values.
type Option func(*client) error

func WithCallOptions(v ...grpc.CallOption) Option {
	return func(o *client) error { o.callOptions = v; return nil }
}

func WithContext(v context.Context) Option {
	return func(o *client) error { o.ctx = v; return nil }
}

func WithDialTimeout(v time.Duration) Option {
	return func(o *client) error { o.dialTimeout = v; return nil }
}

func WithFailFast(v bool) Option {
	return func(o *client) error { o.failFast = v; return nil }
}

func WithID(v []byte) Option {
	return func(o *client) error { o.id = v; return nil }
}

func WithIdleStrategyFast(v cc.IdleStrategy) Option {
	return func(o *client) error { o.idleStrategy = v; return nil }
}

func WithKeepAlivePeriod(v time.Duration) Option {
	return func(o *client) error { o.keepAlivePeriod = v; return nil }
}

func WithLogLevel(v string) Option {
	return func(_ *client) error { logging.SetLevel(v); return nil }
}

func WithLoggerName(v string) Option {
	return func(o *client) error { o.logger = logging.NewLogger(v); return nil }
}

func WithPeers(v string) Option {
	return func(o *client) error { o.peers = v; return nil }
}

func WithRetryTimeout(v time.Duration) Option {
	return func(o *client) error { o.retryTimeout = v; return nil }
}

func WithDNS(host string, expectedSize int, timeout time.Duration) Option {
	return func(o *client) error {
		if expectedSize < 3 || expectedSize%2 == 0 {
			return fmt.Errorf("incorrect cluster size: %d", expectedSize)
		}

		t := time.After(timeout)
		for {
			select {
			case <-t:
				return ErrDeadlineExceeded

			default:
				ips, err := net.LookupIP(host)
				if err != nil {
					return err
				}
				sz := len(ips)
				if sz != expectedSize {
					break
				}

				sort.Slice(ips, func(i, j int) bool { return bytes.Compare(ips[i], ips[j]) < 0 })
				b := strings.Builder{}
				for _, ip := range ips {
					b.WriteString(ip.To4().String())
					b.WriteString(peerListSep)
				}
				o.peers = strings.TrimSuffix(b.String(), peerListSep)
				return nil

			}
		}

	}
}

type State int

type StateChan = <-chan State

const (
	Connecting State = iota + 100
	Connected
)

func (s State) String() string {
	switch s {
	case Connecting:
		return "CONNECTED"
	case Connected:
		return "CONNECTED"
	default:
		return "UNKNOWN"
	}
}

// Client defines Stoa client methods.
type Client interface {
	// Queue returns a Queue with the name n.
	Queue(n string) Queue

	// Dictionary returns a Dictionary with the name n.
	Dictionary(n string) Dictionary

	// Mutex returns a Mutex with the name n.
	Mutex(n string) Mutex

	// State returns a channel which reports the Client state changes.
	State() StateChan
}

// CallOption adds an option to a method call.
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

func (c *client) createBase(name string) base {
	return base{
		name:         name,
		handle:       c.handle,
		callOptions:  c.callOptions,
		failFast:     c.failFast,
		idleStrategy: c.idleStrategy,
		retryTimeout: c.retryTimeout,
	}
}

// Name implements Collection.Name.
func (b *base) Name() string {
	return b.name
}
