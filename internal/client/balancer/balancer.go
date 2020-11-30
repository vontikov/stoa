package balancer

import (
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vontikov/stoa/internal/logging"
)

// Name is the balancer name.
const Name = "stoa_picker"

func init() {
	balancer.Register(newBuilder())
}

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &builder{}, base.Config{HealthCheck: true})
}

type builder struct{}

func (*builder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn
	for sc := range info.ReadySCs {
		scs = append(scs, sc)
	}
	return &picker{
		logger: logging.NewLogger("picker"),
		scs:    scs,
		rcs:    info.ReadySCs,
	}
}

type picker struct {
	sync.Mutex
	logger logging.Logger
	scs    []balancer.SubConn
	rcs    map[balancer.SubConn]base.SubConnInfo
	idx    int
}

func (p *picker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	p.Lock()
	defer p.Unlock()

	return balancer.PickResult{
		SubConn: p.scs[p.idx],
		Done: func(info balancer.DoneInfo) {
			if info.Err == nil {
				return
			}

			// a follower was called
			if s, ok := status.FromError(info.Err); ok && s.Code() == codes.FailedPrecondition {
				p.Lock()
				defer p.Unlock()
				p.idx = (p.idx + 1) % len(p.scs)
				if p.logger.IsTrace() {
					p.logger.Trace("switched to", "address", p.rcs[p.scs[p.idx]].Address.Addr)
				}
			}
		},
	}, nil
}
