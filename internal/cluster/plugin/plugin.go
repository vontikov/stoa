package plugin

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver/dns"

	"github.com/vontikov/stoa/internal/util"
	"github.com/vontikov/stoa/pkg/pb"
)

type FSM interface {
	DictionaryScan() chan *pb.DictionaryEntry
	DictionaryRemove(*pb.DictionaryEntry) error
}

type Plugin struct {
	sync.RWMutex
	target   string
	dialOpts []grpc.DialOption
	conn     *grpc.ClientConn
	handle   pb.EvictionClient
	fsm      FSM
}

var tickerPeriod = 1000 * time.Millisecond
var errHandler util.ErrorHandler

func ErrorHandler(h util.ErrorHandler) { errHandler = h }
func TickerPeriod(d time.Duration)     { tickerPeriod = d }

func New(target string, supplier FSM) *Plugin {
	return &Plugin{
		target: target,
		dialOpts: []grpc.DialOption{
			grpc.WithResolvers(dns.NewBuilder()),
			grpc.WithBlock(),
			grpc.WithInsecure(),
			grpc.WithKeepaliveParams(
				keepalive.ClientParameters{
					Time:                10 * time.Second,
					Timeout:             time.Second,
					PermitWithoutStream: true,
				}),
		},
		fsm: supplier,
	}
}

func (p *Plugin) Ready() bool {
	p.RLock()
	defer p.RUnlock()
	return p.conn != nil && p.conn.GetState() == connectivity.Ready
}

func (p *Plugin) Run(ctx context.Context) error {
	conn, err := grpc.DialContext(ctx, p.target, p.dialOpts...)
	if err != nil {
		return err
	}

	p.Lock()
	p.conn = conn
	p.handle = pb.NewEvictionClient(conn)
	p.Unlock()

	t := time.NewTicker(tickerPeriod)
	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return ctx.Err()
		case <-t.C:
			if err := p.scan(ctx); err != nil && errHandler != nil && !errHandler(err) {
				return err
			}
		}
	}
}

func (p *Plugin) scan(ctx context.Context) error {
	if !p.Ready() {
		return nil
	}

	c := make(chan error, 1)
	go func() {
		for e := range p.fsm.DictionaryScan() {
			r, err := p.handle.Dictionary(ctx, e)
			if err != nil {
				c <- err
				return
			}
			if r.Ok {
				err := p.fsm.DictionaryRemove(e)
				if err != nil {
					c <- err
					return
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c:
		return err
	}
}
