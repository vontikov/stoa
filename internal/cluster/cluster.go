package cluster

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
)

const (
	defaultBindPort    = 3499
	defaultBindAddress = "0.0.0.0"
	defaultLoggerName  = "stoa-raft"
	defaultPingPeriod  = 1000 * time.Millisecond
)

// Cluster is a higher level representation of raft.Raft cluster endpoint.
type Cluster interface {
	// Raft returns the raft node for this Cluster.
	Raft() *raft.Raft

	// ID returns the Cluster ID
	ID() string

	// IsLeader returns true if the Cluster is the leader.
	IsLeader() bool

	// LeaderAddress returns address of the current leader of the cluster.
	LeaderAddress() string

	// AddVoter adds a new Voter to the cluster.
	AddVoter(id string, bindIP string, bindPort int) error

	// LeadershipTransfer will transfer leadership to another Cluster.
	LeadershipTransfer() error

	// Leader returns a channel which signals on acquiring or losing
	// leadership.  It sends true if the Cluster become the leader, otherwise it
	// returns false.
	Leader() <-chan bool

	// Execute executes the command c bypassing Raft node.
	Execute(c *pb.ClusterCommand) interface{}

	// Done returns a channel that's closed when the Cluster is shutdown.
	Done() <-chan error
}

// ErrDeadlineExceeded is the error returned if a timeout is specified.
var ErrDeadlineExceeded error = errors.New("deadline exceeded")

// ErrClusterConfig is the error returned in case of invalid Cluster
// configuration.
var ErrClusterConfig = errors.New("invalid cluster configuration")

// ErrPeerParams is the error returned if the peer parameters are invalid.
var ErrPeerParams = errors.New("invalid peer parameters")

// cluster implements Cluster.
type cluster struct {
	logger logging.Logger
	r      *raft.Raft
	id     string
	f      *FSM

	doneMutex sync.Mutex // protects following fields
	doneChan  chan error

	leaderMutex sync.RWMutex // protects following fields
	leaderChan  chan bool
}

// PeerListSep separates peer definitions in the peer list.
const PeerListSep = ","

// PeerOptsSep separates a peer's options.
const PeerOptsSep = ":"

type peer struct {
	id   string
	addr string
	port int
}

func newPeer(args []string) (*peer, error) {
	var (
		addr string
		port int
		err  error
	)

	switch len(args) {
	case 1:
		addr = args[0]
		port = defaultBindPort
	case 2:
		addr = args[0]
		port, err = strconv.Atoi(args[1])
		if err != nil {
			return nil, ErrPeerParams
		}
	default:
		return nil, ErrPeerParams
	}

	id := peerID(addr, port)
	return &peer{id, addr, port}, nil
}

func parsePeers(in string) ([]*peer, error) {
	var peers []*peer

	s := strings.TrimSpace(in)
	if strings.HasSuffix(s, PeerListSep) {
		s = s[:len(s)-len(PeerListSep)]
	}

	defs := strings.Split(s, PeerListSep)
	for _, d := range defs {
		args := strings.Split(strings.TrimSpace(d), PeerOptsSep)
		p, err := newPeer(args)
		if err != nil {
			return nil, err
		}
		peers = append(peers, p)
	}
	return peers, nil
}

type options struct {
	bindAddr   string
	loggerName string
	peers      []*peer
	pingPeriod time.Duration
}

func newOptions(opts ...Option) (*options, error) {
	cfg := &options{
		bindAddr:   defaultBindAddress,
		loggerName: defaultLoggerName,
		pingPeriod: defaultPingPeriod,
	}

	for _, o := range opts {
		if err := o(cfg); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

// Option defines a Cluster configuration option.
type Option func(*options) error

// WithPeers passes the list of peers with which the Clustir should be created.
func WithPeers(v string) Option {
	return func(o *options) error {
		peers, err := parsePeers(v)
		if err != nil {
			return err
		}
		o.peers = peers
		return nil
	}
}

// WithBindAddress passes the Cluster's bind address.
func WithBindAddress(v string) Option {
	return func(o *options) error {
		o.bindAddr = v
		return nil
	}
}

// WithLoggerName sets the Cluster logger name.
func WithLoggerName(v string) Option {
	return func(o *options) error {
		o.loggerName = v
		return nil
	}
}

// WithPingPeriod sets the Cluster ping period.
func WithPingPeriod(v time.Duration) Option {
	return func(o *options) error {
		o.pingPeriod = v
		return nil
	}
}

// New creates a new Cluster instance.
func New(ctx context.Context, opts ...Option) (c Cluster, err error) {
	cfg, err := newOptions(opts...)
	if err != nil {
		return nil, err
	}

	logger := logging.NewLogger(cfg.loggerName)

	p := cfg.peers[0]
	bindAddr := fmt.Sprintf("%s:%d", cfg.bindAddr, p.port)
	advAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", p.addr, p.port))
	if err != nil {
		return
	}

	logger.Debug("creating transport", "bind address", bindAddr, "advertise address", advAddr)
	trans, err := raft.NewTCPTransport(bindAddr, advAddr, 7, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}

	conf := raft.DefaultConfig()
	conf.LocalID = raft.ServerID(p.id)
	conf.Logger = logger

	logs := raft.NewInmemStore()
	stable := raft.NewInmemStore()
	snaps := raft.NewInmemSnapshotStore()

	logger.Debug("creating raft")
	fsm := NewFSM(ctx)
	r, err := raft.NewRaft(conf, fsm, logs, stable, snaps, trans)
	if err != nil {
		return nil, err
	}

	cl := &cluster{
		logger: logger,
		r:      r,
		id:     p.id,
		f:      fsm,
	}

	if len(cfg.peers) > 1 {
		logger.Debug("bootstrap cluster")
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					Address: raft.ServerAddress(bindAddr),
					ID:      raft.ServerID(p.id),
				},
			},
		}
		if err = r.BootstrapCluster(configuration).Error(); err != nil {
			logger.Error("bootstrap error", "message", err)
			return nil, err
		}
		if err := waitLeaderStatus(r, 10*time.Second); err == nil {
			for _, p := range cfg.peers[1:] {
				addr := raft.ServerAddress(fmt.Sprintf("%s:%d", p.addr, p.port))
				logger.Info("adding voter", "id", p.id, "address", addr)
				if err = r.AddVoter(raft.ServerID(p.id), addr, 0, 0).Error(); err != nil {
					logger.Warn("voter error", "message", err)
				}
			}
		} else {
			logger.Warn("not leader")
		}
	}

	go cl.watchLeadership(ctx)
	go cl.watchExpiration(ctx, cfg)
	go cl.cleanup(ctx)

	c = cl
	return
}

func (c *cluster) watchLeadership(ctx context.Context) {
	logger := c.logger

	logger.Debug("leadership watcher started")
	defer logger.Debug("leadership watcher stopped")

	ch := c.r.LeaderCh()
	for {
		select {
		case <-ctx.Done():
			return
		case leader := <-ch:
			c.f.leader(leader)

			logger.Debug("leadership signal", "status", leader)
			c.leaderMutex.RLock()
			if c.leaderChan != nil {
				c.leaderChan <- leader
			}
			c.leaderMutex.RUnlock()
		}
	}
}

func (c *cluster) watchExpiration(ctx context.Context, cfg *options) {
	logger := c.logger
	pingPeriod := cfg.pingPeriod
	expirationPeriod := pingPeriod * 3

	logger.Debug("expiration watcher started", "ping", pingPeriod, "period", expirationPeriod)
	defer logger.Debug("expiration watcher stopped")

	t := time.NewTicker(cfg.pingPeriod)
	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			_ = mutexUnlockExpired(c.f, expirationPeriod)
		}
	}
}

func (c *cluster) cleanup(ctx context.Context) {
	<-ctx.Done()
	e := c.r.Shutdown().Error()
	c.doneMutex.Lock()
	if c.doneChan != nil {
		c.doneChan <- e
		close(c.doneChan)
	}
	c.doneMutex.Unlock()
	c.logger.Info("shutdown")
}

// Raft implements Cluster.Raft.
func (c *cluster) Raft() *raft.Raft {
	return c.r
}

// ID implements Cluster.ID.
func (c *cluster) ID() string {
	return c.id
}

// LeaderAddress implements Cluster.LeaderAddress.
func (c *cluster) LeaderAddress() string {
	return string(c.r.Leader())
}

// IsLeader implements Cluster.IsLeader.
func (c *cluster) IsLeader() bool {
	return raft.Leader == c.r.State()
}

// AddVoter implements Cluster.AddVoter.
func (c *cluster) AddVoter(id string, bindIP string, bindPort int) error {
	addr := raft.ServerAddress(fmt.Sprintf("%s:%d", bindIP, bindPort))
	c.logger.Info("adding voter", "id", id, "address", addr)
	return c.r.AddVoter(raft.ServerID(id), addr, 0, 0).Error()
}

// LeadershipTransfer implements Cluster.LeadershipTransfer.
func (c *cluster) LeadershipTransfer() error {
	return c.r.LeadershipTransfer().Error()
}

// Leader implements Cluster.Leader.
func (c *cluster) Leader() <-chan bool {
	// TODO
	c.leaderMutex.Lock()
	if c.leaderChan == nil {
		c.leaderChan = make(chan bool, 1)
	}
	ch := c.leaderChan
	c.leaderMutex.Unlock()
	return ch
}

// Execute implements Cluster.Execute.
func (c *cluster) Execute(cmd *pb.ClusterCommand) interface{} {
	return c.f.Execute(cmd)
}

// Done implements Cluster.Done.
func (c *cluster) Done() <-chan error {
	c.doneMutex.Lock()
	defer c.doneMutex.Unlock()
	if c.doneChan == nil {
		c.doneChan = make(chan error, 1)
	}
	return c.doneChan
}

func peerID(bindAddr string, bindPort int) string {
	return fmt.Sprintf("%s:%d", bindAddr, bindPort)
}

func waitLeaderStatus(r *raft.Raft, d time.Duration) error {
	const sleepDuration = 100 * time.Millisecond
	timeout := time.After(d)
	for {
		select {
		case <-timeout:
			return ErrDeadlineExceeded
		default:
			if raft.Leader == r.State() {
				return nil
			}
			time.Sleep(sleepDuration)
		}
	}
}
