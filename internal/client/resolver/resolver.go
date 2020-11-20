package resolver

import (
	"fmt"

	"google.golang.org/grpc/resolver"

	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/logging"
)

const (
	// Scheme contains the resolver scheme
	Scheme = "raft"
)

// StaticBuilder implements
// ResolverBuilder(https://godoc.org/google.golang.org/grpc/resolver#Builder).
type StaticBuilder struct {
	peers []*cluster.Peer
}

// staticResolver implements
// Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
type staticResolver struct {
	logger logging.Logger
	cc     resolver.ClientConn
	target resolver.Target
	peers  []*cluster.Peer
}

// New creates new resolver.Builder instance
func New(peers string) (resolver.Builder, error) {
	peerList, err := cluster.ParsePeers(peers)
	if err != nil {
		return nil, err
	}
	return &StaticBuilder{peers: peerList}, nil
}

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
		addr := fmt.Sprintf("%s:%d", p.BindAddr, p.BindPort)
		addresses = append(addresses, resolver.Address{Addr: addr})
	}

	r.cc.UpdateState(resolver.State{Addresses: addresses})
	r.logger.Debug("updated with", "addresses", addresses)
	return
}
