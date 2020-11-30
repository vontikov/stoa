package resolver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/grpc/resolver"

	"github.com/vontikov/stoa/internal/logging"
)

const (
	// Scheme contains the resolver scheme.
	Scheme = "raft"

	// PeerListSep separates peer definitions in the peer list.
	PeerListSep = ","

	// PeerOptsSep separates a peer's options.
	PeerOptsSep = ":"

	// DefaultGRPCPort is the default gRPC port.
	DefaultGRPCPort = 3500
)

// ErrPeerParams is the error returned if the peer parameters are invalid.
var ErrPeerParams = errors.New("invalid peer parameters")

// ErrEmptyPeerList is the error returned if an empty peer list passed to the
// Cluster configuration.
var ErrEmptyPeerList = errors.New("empty peer list")

type peer struct {
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
		port = DefaultGRPCPort
	case 2:
		addr = args[0]
		port, err = strconv.Atoi(args[1])
		if err != nil {
			return nil, ErrPeerParams
		}
	default:
		return nil, ErrPeerParams
	}
	return &peer{addr, port}, nil
}

func parsePeers(in string) ([]*peer, error) {
	var peers []*peer

	s := strings.TrimSpace(in)
	if s == "" {
		return nil, ErrEmptyPeerList
	}
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

// StaticBuilder implements
// ResolverBuilder(https://godoc.org/google.golang.org/grpc/resolver#Builder).
type StaticBuilder struct {
	peers []*peer
}

// staticResolver implements
// Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
type staticResolver struct {
	logger logging.Logger
	cc     resolver.ClientConn
	target resolver.Target
	peers  []*peer
}

// New creates new resolver.Builder instance
func New(peers string) (resolver.Builder, error) {
	peerList, err := parsePeers(peers)
	if err != nil {
		return nil, err
	}
	return &StaticBuilder{peers: peerList}, nil
}

// Build implements resolver.Builder.Build.
func (b *StaticBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := staticResolver{
		logger: logging.NewLogger("static-resolver"),
		cc:     cc,
		target: target,
		peers:  b.peers,
	}
	r.resolve()
	return &r, nil
}

// Scheme implements resolver.Builder.Scheme.
func (*StaticBuilder) Scheme() string {
	return Scheme
}

func (r *staticResolver) ResolveNow(resolver.ResolveNowOptions) {
	r.resolve()
}

func (r *staticResolver) Close() {}

func (r *staticResolver) resolve() {
	var addresses []resolver.Address
	for _, p := range r.peers {
		addr := fmt.Sprintf("%s:%d", p.addr, p.port)
		addresses = append(addresses, resolver.Address{Addr: addr})
	}

	r.cc.UpdateState(resolver.State{Addresses: addresses})
	if r.logger.IsTrace() {
		r.logger.Trace("updated with", "addresses", addresses)
	}
	return
}
