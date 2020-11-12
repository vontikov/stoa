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

func queue(f *FSM, n string) cc.Queue {
	v, _ := f.qs.ComputeIfAbsent(n, func() interface{} { return newQueueRecord() })
	return v.(*queueRecord).q
}

func queueSize(f *FSM, m *pb.ClusterCommand) interface{} {
	q := queue(f, m.GetName().Name)
	v := make([]byte, 4)
	binary.LittleEndian.PutUint32(v, uint32(q.Size()))
	return &pb.Value{Value: v}
}

func queueClear(f *FSM, m *pb.ClusterCommand) interface{} {
	q := queue(f, m.GetName().Name)
	q.Clear()
	return nil
}

func queueOffer(f *FSM, m *pb.ClusterCommand) interface{} {
	v := m.GetValue()
	q := queue(f, v.Name)
	e := entry{
		Value:     v.Value,
		TTLMillis: m.TtlMillis,
	}
	q.Offer(e)
	return nil
}

func queuePoll(f *FSM, m *pb.ClusterCommand) interface{} {
	q := queue(f, m.GetName().Name)
	if v := q.Poll(); v != nil {
		return &pb.Value{Value: v.(entry).Value}
	}
	return &pb.Value{}
}

func queuePeek(f *FSM, m *pb.ClusterCommand) interface{} {
	q := queue(f, m.GetName().Name)
	if v := q.Peek(); v != nil {
		return &pb.Value{Value: v.(entry).Value}
	}
	return &pb.Value{}
}
