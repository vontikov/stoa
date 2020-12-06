package client

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	cc "github.com/vontikov/go-concurrent"

	"github.com/vontikov/stoa/internal/client/balancer"
	"github.com/vontikov/stoa/internal/client/resolver"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
)

// ErrNotReady is the error returned by client.Ready when the client is not
// ready.
var ErrNotReady = errors.New("client not ready")

// client is the Client implementation.
type client struct {
	sync.RWMutex
	sync.WaitGroup
	logger logging.Logger
	cfg    *options
	conn   *grpc.ClientConn
	handle pb.StoaClient
	qs     cc.Map
	ds     cc.Map
	ms     cc.Map
}

// New creates and returns a new Client instance.
func New(opts ...Option) (Client, error) {
	cfg := newOptions()
	for _, o := range opts {
		o(cfg)
	}
	cfg.applyDefaults()

	resolver, err := resolver.New(cfg.peers)
	if err != nil {
		return nil, err
	}
	dialOpts := []grpc.DialOption{
		grpc.WithResolvers(resolver),
		grpc.WithBalancerName(balancer.Name),
		//		grpc.WithBalancerName("pick_first"),
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                10 * time.Second,
				Timeout:             time.Second,
				PermitWithoutStream: true,
			}),
	}

	c := client{
		logger: logging.NewLogger("stoa"),
		cfg:    cfg,
		qs:     cc.NewSynchronizedMap(0),
		ds:     cc.NewSynchronizedMap(0),
		ms:     cc.NewSynchronizedMap(0),
	}
	c.Add(1)
	go func() {
		defer c.Done()
		c.dial(cfg, dialOpts)
	}()

	if err := c.ready(cfg.dialTimeout); err != nil {
		return nil, err
	}

	/*
		if err := c.watcher(); err != nil {
			return nil, err
		}
	*/
	return &c, nil
}

func (c *client) dial(cfg *options, dialOpts []grpc.DialOption) error {
	target := fmt.Sprintf("%s:///%s", resolver.Scheme, "stoa")
	conn, err := grpc.DialContext(cfg.context, target, dialOpts...)
	if err != nil {
		return err
	}

	c.Lock()
	handle := pb.NewStoaClient(conn)
	c.conn = conn
	c.handle = handle
	c.Unlock()
	<-cfg.context.Done()
	return conn.Close()
}

func (c *client) ready(d time.Duration) error {
	t := time.After(d)
	idle := cc.NewSleepingIdleStrategy(100 * time.Millisecond)

	for {
		select {
		case <-t:
			return ErrNotReady
		default:
			if func() bool {
				c.RLock()
				defer c.RUnlock()
				return c.conn != nil && c.conn.GetState() == connectivity.Ready
			}() {
				return nil
			}
			idle.Idle()
		}
	}
}

func (c *client) watcher() error {
	stream, err := c.handle.Keep(c.cfg.context)
	if err != nil {
		c.logger.Error("stream error", "message", err)
		return err
	}
	ctx := c.cfg.context

	// ping
	c.Add(1)
	go func() {
		defer c.Done()
		t := time.NewTicker(c.cfg.keepAlivePeriod)
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if c.logger.IsTrace() {
					c.logger.Trace("sending ping", "id", c.cfg.id)
				}
				if err := stream.Send(&pb.Ping{Id: c.cfg.id}); err != nil {
					c.logger.Warn("ping error", "message", err)
				}
			}
		}
	}()

	// status
	c.Add(1)
	go func() {
		defer c.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				statusMsg, err := stream.Recv()
				if err != nil {
					s, ok := status.FromError(err)
					if ok {
						if s.Code() == codes.Canceled {
							return
						}
					}
					if s.Code() != codes.FailedPrecondition {
						c.logger.Warn("status error", "code", s.Code(), "message", err)
					}
				}
				if c.logger.IsTrace() {
					c.logger.Trace("status received", "status", statusMsg)
				}

				switch statusMsg.GetPayload().(type) {
				case *pb.Status_Mutex:
					mx := statusMsg.GetMutex()
					if v := c.ms.Get(mx.Name); v != nil {
						v.(*mutex).watchers.Range(
							func(e interface{}) bool {
								e.(MutexWatcher).Apply(mx.Name, mx.Locked)
								return true
							})
					}
				}
			}
		}
	}()

	c.logger.Info("watcher started")
	return nil
}

func (c *client) Queue(n string) Queue {
	v, _ := c.qs.ComputeIfAbsent(n, func() interface{} {
		return newQueue(n, c.cfg, c.handle)
	})
	return v.(Queue)
}

func (c *client) Dictionary(n string) Dictionary {
	v, _ := c.ds.ComputeIfAbsent(n, func() interface{} {
		return newDictionary(n, c.cfg, c.handle)
	})
	return v.(Dictionary)
}

func (c *client) Mutex(n string) Mutex {
	v, _ := c.ms.ComputeIfAbsent(n, func() interface{} {
		return newMutex(n, c.cfg, c.handle)
	})
	return v.(Mutex)
}
