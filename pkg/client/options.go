package client

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/google/uuid"
	cc "github.com/vontikov/go-concurrent"
	"github.com/vontikov/stoa/internal/logging"
)

type options struct {
	callOptions     []grpc.CallOption
	context         context.Context
	dialTimeout     time.Duration
	failFast        bool
	id              string
	idleStrategy    cc.IdleStrategy
	keepAlivePeriod time.Duration
	peers           string
	retryTimeout    time.Duration
}

// Option is a function applied to an options to change the options' default values.
type Option func(*options)

func WithCallOptions(v []grpc.CallOption) Option    { return func(o *options) { o.callOptions = v } }
func WithContext(v context.Context) Option          { return func(o *options) { o.context = v } }
func WithDialTimeout(v time.Duration) Option        { return func(o *options) { o.dialTimeout = v } }
func WithFailFast(v bool) Option                    { return func(o *options) { o.failFast = v } }
func WithID(v string) Option                        { return func(o *options) { o.id = v } }
func WithIdleStrategyFast(v cc.IdleStrategy) Option { return func(o *options) { o.idleStrategy = v } }
func WithKeepAlivePeriod(v time.Duration) Option    { return func(o *options) { o.keepAlivePeriod = v } }
func WithLogLevel(v string) Option                  { return func(_ *options) { logging.SetLevel(v) } }
func WithPeers(v string) Option                     { return func(o *options) { o.peers = v } }
func WithRetryTimeout(v time.Duration) Option       { return func(o *options) { o.retryTimeout = v } }

func newOptions() *options {
	return &options{
		id: uuid.New().String(),
	}
}

func (o *options) applyDefaults() {
	if o.context == nil {
		o.context = context.Background()
	}
	if o.callOptions == nil {
		o.callOptions = []grpc.CallOption{grpc.WaitForReady(true)}
	}
	if o.dialTimeout == 0 {
		o.dialTimeout = DefaultDialTimeout
	}
	if o.retryTimeout == 0 {
		o.retryTimeout = DefaultRetryTimeout
	}
	if o.keepAlivePeriod == 0 {
		o.keepAlivePeriod = DefaultKeepAlivePeriod
	}
	if o.idleStrategy == nil {
		o.idleStrategy = cc.NewSleepingIdleStrategy(500 * time.Millisecond)
	}
}
