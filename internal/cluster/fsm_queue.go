package cluster

import (
	"encoding/binary"

	cc "github.com/vontikov/go-concurrent"
	"github.com/vontikov/stoa/pkg/pb"
)

var initialQueueCapacity = 16

type queueRecord struct {
	q cc.Queue
}

func newQueueRecord() *queueRecord {
	return &queueRecord{q: cc.NewSynchronizedRingQueue(initialQueueCapacity)}
}

func (f *FSM) queue(n string) cc.Queue {
	v, _ := f.qs.ComputeIfAbsent(n, func() interface{} { return newQueueRecord() })
	return v.(*queueRecord).q
}

func (f *FSM) queueSize(m *pb.ClusterCommand) interface{} {
	q := f.queue(m.GetName().Name)
	v := make([]byte, 4)
	binary.LittleEndian.PutUint32(v, uint32(q.Size()))
	return &pb.Value{Value: v}
}

func (f *FSM) queueClear(m *pb.ClusterCommand) interface{} {
	q := f.queue(m.GetName().Name)
	q.Clear()
	return nil
}

func (f *FSM) queueOffer(m *pb.ClusterCommand) interface{} {
	v := m.GetValue()
	q := f.queue(v.Name)
	e := entry{
		Value:     v.Value,
		TTLMillis: m.TtlMillis,
	}
	q.Offer(e)
	return nil
}

func (f *FSM) queuePoll(m *pb.ClusterCommand) interface{} {
	q := f.queue(m.GetName().Name)
	if v := q.Poll(); v != nil {
		return &pb.Value{Value: v.(entry).Value}
	}
	return &pb.Value{}
}

func (f *FSM) queuePeek(m *pb.ClusterCommand) interface{} {
	q := f.queue(m.GetName().Name)
	if v := q.Peek(); v != nil {
		return &pb.Value{Value: v.(entry).Value}
	}
	return &pb.Value{}
}
