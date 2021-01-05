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

const defaultBindPort = 3499
const defaultBindAddress = "0.0.0.0"

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

	// WatchLeadership returns a channel which signals on acquiring or losing
	// leadership.  It sends true if the Cluster become the leader, otherwise it
	// returns false.
	WatchLeadership() <-chan bool

	// Apply executes the command c bypassing Raft node.
	Apply(c *pb.ClusterCommand) interface{}

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

	mu   sync.Mutex
	done chan error
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
	autoDiscovery bool
	bindAddr      string
	peers         []*peer
}

func newOptions(opts ...Option) (*options, error) {
	cfg := &options{}
	for _, o := range opts {
		if err := o(cfg); err != nil {
			return nil, err
		}
	}
	if !cfg.autoDiscovery && len(cfg.peers) == 0 {
		return nil, ErrClusterConfig
	}
	if cfg.bindAddr == "" {
		cfg.bindAddr = defaultBindAddress
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

// New creates a new Cluster instance.
func New(ctx context.Context, opts ...Option) (c Cluster, err error) {
	cfg, err := newOptions(opts...)
	if err != nil {
		return nil, err
	}

	logger := logging.NewLogger("raft")

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

	cl := &cluster{
		logger: logger,
		r:      r,
		id:     p.id,
		f:      fsm,
	}

	go func() {
		<-ctx.Done()

		cl.mu.Lock()
		if cl.done == nil {
			cl.done = make(chan error, 1)
		}
		d := cl.done
		cl.mu.Unlock()

		d <- r.Shutdown().Error()
		close(d)
		logger.Info("shutdown")
	}()

	c = cl
	return
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

// WatchLeadership implements Cluster.WatchLeadership.
func (c *cluster) WatchLeadership() <-chan bool {
	return c.r.LeaderCh()
}

// Apply implements Cluster.Apply.
func (c *cluster) Apply(cmd *pb.ClusterCommand) interface{} {
	return c.f.Execute(cmd)
}

// Done implements Cluster.Done.
func (c *cluster) Done() <-chan error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.done == nil {
		c.done = make(chan error, 1)
	}
	return c.done
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
