package resolver

import (
	"testing"

	"context"

	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"

	"github.com/stretchr/testify/assert"
)

type testSupplier struct {
	peers []peer
}

type testConn struct {
	state resolver.State
	err   error
}

func newTestSupplier() *testSupplier { return &testSupplier{} }

func (s *testSupplier) add(p peer) { s.peers = append(s.peers, p) }
func (s *testSupplier) Get(string) <-chan peer {
	ch := make(chan peer)
	go func() {
		for _, p := range s.peers {
			ch <- p
		}
		close(ch)
	}()
	return ch
}
func (*testSupplier) Run(context.Context) error { return nil }

func (*testConn) NewAddress([]resolver.Address)                        {}
func (*testConn) NewServiceConfig(string)                              {}
func (*testConn) ParseServiceConfig(string) *serviceconfig.ParseResult { return nil }

func (t *testConn) ReportError(e error) {
	t.err = e
}

func (t *testConn) UpdateState(s resolver.State) {
	t.state = s
}

func TestResolverTimeoutNoPeer(t *testing.T) {
	s := newTestSupplier()
	b := New(s)
	r, _ := b.Build(resolver.Target{}, &testConn{}, resolver.BuildOptions{})
	assert.Equal(t, ErrLeaderNotFound, r.(*raftResolver).cc.(*testConn).err)
}

func TestResolverTimeoutNoLeader(t *testing.T) {
	s := newTestSupplier()
	s.add(peer{ID: "id"})
	b := New(s)
	r, _ := b.Build(resolver.Target{}, &testConn{}, resolver.BuildOptions{})
	assert.Equal(t, ErrLeaderNotFound, r.(*raftResolver).cc.(*testConn).err)
}

func TestResolverSuccess(t *testing.T) {
	s := newTestSupplier()
	s.add(peer{ID: "id", Leader: true})
	b := New(s)
	r, _ := b.Build(resolver.Target{}, &testConn{}, resolver.BuildOptions{})
	assert.Equal(t, nil, r.(*raftResolver).cc.(*testConn).err)
}
