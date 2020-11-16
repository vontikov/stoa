package cluster

import (
	"github.com/vontikov/stoa/pkg/pb"
)

func (f *FSM) processPing(p *pb.Ping) {
	f.logger.Debug("Ping received", "message", p)
	keys := f.ms.Keys()
	for _, k := range keys {
		v := f.ms.Get(k)
		if v == nil {
			continue
		}
		mx := v.(*mutexRecord)
		mx.touch(p.Id)
	}
}
