package cluster

import (
	"github.com/vontikov/stoa/pkg/pb"
)

func processPing(f *FSM, m *pb.ClusterCommand) interface{} {
	id := m.GetClientId().Id

	if f.logger.IsTrace() {
		f.logger.Trace("ping received", "id", string(id))
	}

	for _, k := range f.ms.Keys() {
		if v := f.ms.Get(k); v != nil {
			v.(*mutexRecord).touch(id)
		}
	}
	return nil
}
