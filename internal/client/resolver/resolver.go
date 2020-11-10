package resolver

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/resolver"

	"github.com/vontikov/stoa/internal/logging"
)

const (
	// Scheme contains resolver scheme
	Scheme = "raft"
	// DefaultResolutionRate contains the watcher's default resolution rate
	DefaultResolutionRate = 500 * time.Millisecond
)

// ErrLeaderNotFound is the error reported when there is no Raft lider yet.
var ErrLeaderNotFound = errors.New("leader not found")

var resRate = DefaultResolutionRate

var resolveFuncRef atomic.Value

// builder implements
// ResolverBuilder(https://godoc.org/google.golang.org/grpc/resolver#Builder).
type builder struct {
	supplier Supplier
}

// raftResolver implements
// Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
type raftResolver struct {
	sync.Mutex
	logger   logging.Logger
	cancel   context.CancelFunc
	cc       resolver.ClientConn
	ctx      context.Context
	leaderID string
	supplier Supplier
	target   resolver.Target
	wg       sync.WaitGroup
}

// Rate sets watcher resolution rate
func Rate(d time.Duration) {
	resRate = d
}

// Resolve forces name resolution
func Resolve() {
	resolveFuncRef.Load().(func())()
}

// New creates new resolver.Builder instance
func New(s Supplier) resolver.Builder {
	return &builder{
		supplier: s,
	}
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ctx, cancel := context.WithCancel(context.Background())
	r := &raftResolver{
		logger:   logging.NewLogger("resolver"),
		target:   target,
		cc:       cc,
		ctx:      ctx,
		cancel:   cancel,
		supplier: b.supplier,
	}

	resolveFuncRef.Store(r.resolve)

	r.resolve()
	r.wg.Add(1)
	go r.watcher()
	return r, nil
}

func (*builder) Scheme() string {
	return Scheme
}

func (r *raftResolver) ResolveNow(o resolver.ResolveNowOptions) {
	r.resolve()
}

func (r *raftResolver) Close() {
	r.Lock()
	defer r.Unlock()
	r.cancel()
	r.wg.Wait()
}

func (r *raftResolver) resolve() {
	r.Lock()
	defer r.Unlock()
	for p := range r.supplier.Get(r.target.Endpoint) {
		if p.Leader && p.ID != r.leaderID {
			r.leaderID = p.ID
			addr := fmt.Sprintf("%s:%d", p.GrpcIP, p.GrpcPort)
			r.cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: addr}}})
			r.logger.Info("updated", "address", addr)
			return
		}
	}
	r.cc.ReportError(ErrLeaderNotFound)
}

func (r *raftResolver) watcher() {
	defer r.wg.Done()
	t := time.NewTicker(resRate)
	for {
		select {
		case <-r.ctx.Done():
			t.Stop()
			return
		case <-t.C:
			r.resolve()
		}
	}
}
