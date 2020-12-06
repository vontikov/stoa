package cluster

import (
	"github.com/vontikov/stoa/pkg/pb"
)

func (f *FSM) processPing(p *pb.Ping) {
	if f.logger.IsTrace() {
		f.logger.Trace("ping received", "id", p.Id)
	}
	for _, k := range f.ms.Keys() {
		if v := f.ms.Get(k); v != nil {
			v.(*mutexRecord).touch(p.Id)
		}
	}
}
