package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	cc "github.com/vontikov/go-concurrent"

	"github.com/vontikov/stoa/internal/client/balancer"
	"github.com/vontikov/stoa/internal/client/resolver"
	"github.com/vontikov/stoa/internal/common"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
)

// ErrNotReady is the error returned by client.Ready when the client is not
// ready.
var ErrNotReady = errors.New("client not ready")

// ErrWatchHandshake is the error returned if the watcher handshake failed.
var ErrWatchHandshake = errors.New("watch handshake")

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

	if err := c.watch(); err != nil {
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
				c.logger.Warn("ping error", "message", err)
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

	c.qs.Range(func(k, v interface{}) bool {
		q := v.(*queue)
		q.mu.RLock()
		if q.watch != nil {
			close(q.watch)
		}
		q.watch = nil
		q.mu.RUnlock()
		c.logger.Debug("watch channel closed", "queue", q.name)
		return true
	})

	c.mu.Lock()
	_ = c.conn.Close()
	c.logger.Debug("connection closed")
	c.mu.Unlock()
}

func (c *client) watch() (err error) {
	ctx := c.ctx
	connChan := make(chan bool, 10)
	streamChan := make(chan pb.Stoa_WatchClient)

	var once sync.Once
	errChan := make(chan error)

	go func() {
		id := &pb.ClientId{Id: c.id}

		for {
			select {
			case <-ctx.Done():
				return
			case <-connChan:
				c.reportState(Connecting)
				err := retry(
					ctx,
					func() error {
						stream, err := c.handle.Watch(ctx, id)
						if err != nil {
							return err
						}
						if err := c.handshake(stream); err != nil {
							return err
						}
						streamChan <- stream
						return nil
					},
					c.retryTimeout,
					c.idleStrategy,
					func() { c.logger.Debug("watch connection retry") },
				)

				once.Do(func() { errChan <- err })
				if err != nil {
					return
				}
			}
		}
	}()

	go func() {
		for {
			connChan <- true
			select {
			case <-ctx.Done():
				return
			case stream := <-streamChan:
				c.reportState(Connected)
				if err := c.watchLoop(ctx, stream); err == nil {
					return
				}
				c.logger.Debug("force watcher reconnection")
			}
		}
	}()

	return <-errChan
}

func (c *client) handshake(stream pb.Stoa_WatchClient) error {
	st, err := stream.Recv()
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.FailedPrecondition {
			c.logger.Debug("watch handshake failed", "message", err)
			return ErrWatchHandshake
		}

		c.logger.Error("watch handshake error", "message", err)
		return err
	}

	id := st.GetId()
	if id == nil {
		c.logger.Error("watch handshake unexpected type")
		return ErrWatchHandshake
	}
	if !bytes.Equal(id.Id, c.id) {
		c.logger.Error("watch handshake unexpected Client ID", "actual", id.Id)
		return ErrWatchHandshake
	}

	c.logger.Debug("watch handshake success")
	return nil
}

func (c *client) watchLoop(ctx context.Context, stream pb.Stoa_WatchClient) error {
	c.logger.Debug("watch loop started")
	defer c.logger.Debug("watch loop completed")

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			st, err := stream.Recv()
			if err == nil {
				if err := c.processStatus(st); err != nil {
					return err
				}
				break
			}

			if err == io.EOF {
				return nil
			}
			if s, ok := status.FromError(err); ok {
				switch s.Code() {
				case codes.Unavailable, codes.Canceled:
					return nil
				case codes.FailedPrecondition:
					return err
				default:
					c.logger.Debug("status error", "message", err)
					return err
				}
			}
			c.logger.Warn("status error", "message", err)
			return err
		}
	}
}

func (c *client) processStatus(st *pb.Status) error {
	if c.logger.IsTrace() {
		c.logger.Trace("status received", "data", st)
	}

	switch st.U.(type) {
	case *pb.Status_C:
		if st.GetC().S == pb.ClusterStatus_LEADERSHIP_LOST {
			return errors.New("leadership lost")
		}
	case *pb.Status_Q:
		qs := st.GetQ()
		if v := c.qs.Get(qs.Name); v != nil {
			q := v.(*queue)
			q.mu.RLock()
			w := q.watch
			q.mu.RUnlock()
			if w != nil {
				select {
				case w <- qs:
				default:
					c.logger.Warn("queue watch channel blocked", "name", qs.Name)
				}
			}
		}
	case *pb.Status_M:
		ms := st.GetM()
		if v := c.ms.Get(ms.Name); v != nil {
			m := v.(*mutex)
			m.mu.RLock()
			w := m.watch
			m.mu.RUnlock()
			if w != nil {
				select {
				case w <- ms:
				default:
					c.logger.Warn("mutex watch channel blocked", "name", ms.Name)
				}
			}
		}
	}
	return nil
}
