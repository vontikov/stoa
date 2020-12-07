package cluster

import (
	"github.com/vontikov/stoa/pkg/pb"
)

func processPing(f *FSM, m *pb.ClusterCommand) interface{} {
	cid := m.GetCid()
	id := cid.Id
	if f.logger.IsTrace() {
		f.logger.Trace("ping received", "id", id)
	}

	for _, k := range f.ms.Keys() {
		if v := f.ms.Get(k); v != nil {
			v.(*mutexRecord).touch(id)
		}
	}
	return &pb.Empty{}
}
