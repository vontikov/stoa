package cluster

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	"github.com/vontikov/stoa/internal/logging"
)

// Peer is a higher level representation of Raft cluster endpoint.
type Peer struct {
	sync.Mutex
	R            *raft.Raft
	id           string
	bindIP       string
	bindPort     int
	localAddress raft.ServerAddress
}

const network = "tcp"

// NewPeer creates new Peer with the id, listening incoming requests on bindIP and bindPort
func NewPeer(id string, bindIP string, bindPort int, fsm raft.FSM) (*Peer, error) {
	bindAddress := fmt.Sprintf("%s:%d", bindIP, bindPort)
	localAddress := raft.ServerAddress(bindAddress)

	addr, err := net.ResolveTCPAddr(network, bindAddress)
	if err != nil {
		return nil, err
	}
	trans, err := raft.NewTCPTransport(bindAddress, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}

	conf := raft.DefaultConfig()
	conf.LocalID = raft.ServerID(id)
	conf.Logger = logging.NewLogger("raft")

	logs := raft.NewInmemStore()
	stable := raft.NewInmemStore()
	snaps := raft.NewInmemSnapshotStore()

	r, err := raft.NewRaft(conf, fsm, logs, stable, snaps, trans)
	if err != nil {
		return nil, err
	}

	servers := []raft.Server{
		{
			Address: localAddress,
			ID:      raft.ServerID(id),
		},
	}
	configuration := raft.Configuration{
		Servers: servers,
	}
	fut := r.BootstrapCluster(configuration)
	if err := fut.Error(); err != nil {
		return nil, err
	}

	return &Peer{
		R:            r,
		id:           id,
		bindIP:       bindIP,
		bindPort:     bindPort,
		localAddress: localAddress,
	}, nil
}

// ID returns the ID of the Peer p
func (p *Peer) ID() string { return p.id }

// BindIP returns the bind IP of the Peer p
func (p *Peer) BindIP() string { return p.bindIP }

// BindPort returns the bind port of the Peer p
func (p *Peer) BindPort() int { return p.bindPort }

// Shutdown stops the Peer
func (p *Peer) Shutdown() error {
	f := p.R.Shutdown()
	return f.Error()
}

// LeadershipTransfer will transfer leadership to another Peer.
func (p *Peer) LeadershipTransfer() error {
	f := p.R.LeadershipTransfer()
	return f.Error()
}

// AddVoter adds a new Voter to the Peer p
func (p *Peer) AddVoter(id string, bindIP string, bindPort int) error {
	bindAddress := fmt.Sprintf("%s:%d", bindIP, bindPort)
	addr := raft.ServerAddress(bindAddress)
	p.Lock()
	defer p.Unlock()
	return p.R.AddVoter(raft.ServerID(id), addr, 0, 0).Error()
}

// State returns the state of the Peer p
func (p *Peer) State() raft.RaftState {
	return p.R.State()
}

// Leader returns address of the current leader of the cluster
func (p *Peer) Leader() raft.ServerAddress {
	return p.R.Leader()
}

// IsLeader returns true if the Peer p is the leader
func (p *Peer) IsLeader() bool {
	p.Lock()
	defer p.Unlock()
	return raft.Leader == p.R.State()
}
