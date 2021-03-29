package client

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/keepalive"

	cc "github.com/vontikov/go-concurrent"

	"github.com/vontikov/stoa/pkg/client/balancer"
	"github.com/vontikov/stoa/pkg/client/resolver"
	"github.com/vontikov/stoa/pkg/common"
	"github.com/vontikov/stoa/pkg/logging"
	"github.com/vontikov/stoa/pkg/pb"
)

// ErrNotReady is the error returned by client.Ready when the client is not
// ready.
var ErrNotReady = errors.New("client not ready")

// client is the Client implementation.
type client struct {
	callOptions     []grpc.CallOption
	ctx             context.Context
	dialTimeout     time.Duration
	failFast        bool
	id              []byte
	idleStrategy    cc.IdleStrategy
	keepAlivePeriod time.Duration
	logger          logging.Logger
	loggingLevel    string
	loggerName      string
	bootstrap       string
	pingPeriod      time.Duration
	retryTimeout    time.Duration

	qs cc.Map
	ds cc.Map
	ms cc.Map

	mu     sync.RWMutex // protects following fields
	conn   *grpc.ClientConn
	handle pb.StoaClient

	stateMutex sync.RWMutex // protects following fields
	stateChan  chan State
}

// New creates and returns a new Client instance.
func New(ctx context.Context, opts ...Option) (Client, error) {
	c := client{
		ctx:             ctx,
		callOptions:     []grpc.CallOption{grpc.WaitForReady(true)},
		dialTimeout:     defaultDialTimeout,
		id:              genID(common.ClientIDSize),
		idleStrategy:    cc.NewSleepingIdleStrategy(200 * time.Millisecond),
		keepAlivePeriod: defaultKeepAlivePeriod,
		loggingLevel:    defaultLoggerLevel,
		loggerName:      defaultLoggerName,
		pingPeriod:      defaultPingPeriod,
		retryTimeout:    defaultRetryTimeout,

		qs: cc.NewSynchronizedMap(0),
		ds: cc.NewSynchronizedMap(0),
		ms: cc.NewSynchronizedMap(0),
	}

	for _, o := range opts {
		o(&c)
	}

	logging.SetLevel(c.loggingLevel)
	c.logger = logging.NewLogger(c.loggerName)

	if err := c.dial(); err != nil {
		return nil, err
	}

	go c.ping()
	go c.cleanup()

	return &c, nil
}

// Queue implements Client.Queue.
func (c *client) Queue(n string) Queue {
	v, _ := c.qs.ComputeIfAbsent(n, func() interface{} {
		return c.newQueue(n)
	})
	return v.(Queue)
}

// Dictionary implements Client.Dictionary.
func (c *client) Dictionary(n string) Dictionary {
	v, _ := c.ds.ComputeIfAbsent(n, func() interface{} {
		return c.newDictionary(n)
	})
	return v.(Dictionary)
}

// Mutex implements Client.Mutex.
func (c *client) Mutex(n string) Mutex {
	v, _ := c.ms.ComputeIfAbsent(n, func() interface{} {
		return c.newMutex(n)
	})
	return v.(Mutex)
}

// State implements Client.State.
func (c *client) State() <-chan State {
	c.stateMutex.Lock()
	if c.stateChan == nil {
		c.stateChan = make(chan State)
	}
	r := c.stateChan
	c.stateMutex.Unlock()
	return r
}

func (c *client) reportState(s State) {
	if c.logger.IsTrace() {
		c.logger.Trace("report state", "message", s)
	}

	c.stateMutex.RLock()
	defer c.stateMutex.RUnlock()
	if c.stateChan == nil {
		if c.logger.IsTrace() {
			c.logger.Trace("state channel is nil")
		}
		return
	}

	select {
	case c.stateChan <- s:
	default:
		if c.logger.IsTrace() {
			c.logger.Trace("state channel blocked")
		}
	}
}

func (c *client) dial() error {
	g, ctx := errgroup.WithContext(c.ctx)

	target := fmt.Sprintf("%s:///%s", resolver.Scheme, "stoa")

	resolver, err := resolver.New(c.bootstrap)
	if err != nil {
		return err
	}
	dialOpts := []grpc.DialOption{
		grpc.WithResolvers(resolver),
		grpc.WithBalancerName(balancer.Name),
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                10 * time.Second,
				Timeout:             time.Second,
				PermitWithoutStream: true,
			}),
	}

	g.Go(func() error {
		conn, err := grpc.DialContext(ctx, target, dialOpts...)
		if err != nil {
			return err
		}

		c.mu.Lock()
		handle := pb.NewStoaClient(conn)
		c.conn = conn
		c.handle = handle
		c.mu.Unlock()
		return nil
	})

	g.Go(func() error {
		t := time.After(c.dialTimeout)
		for {
			select {
			case <-t:
				return ErrNotReady
			default:
				c.mu.RLock()
				r := c.conn != nil && c.conn.GetState() == connectivity.Ready
				c.mu.RUnlock()
				if r {
					return nil
				}
				c.idleStrategy.Idle()
			}
		}
	})

	return g.Wait()
}

func (c *client) ping() {
	ctx := c.ctx
	t := time.NewTicker(c.pingPeriod)
	m := pb.ClientId{Id: c.id}
	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			if c.logger.IsTrace() {
				c.logger.Trace("sending ping", "id", c.id)
			}
			_, err := c.handle.Ping(ctx, &m, c.callOptions...)
			if err != nil {
				c.logger.Debug("ping error", "message", err)
			}
		}
	}
}

func (c *client) cleanup() {
	<-c.ctx.Done()
	c.logger.Debug("cleaninig up...")

	c.stateMutex.Lock()
	if c.stateChan != nil {
		close(c.stateChan)
	}
	c.stateMutex.Unlock()

	c.mu.Lock()
	_ = c.conn.Close()
	c.logger.Debug("connection closed")
	c.mu.Unlock()
}
