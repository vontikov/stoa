package cluster

import (
	"encoding/binary"

	cc "github.com/vontikov/go-concurrent"

	"github.com/vontikov/stoa/pkg/pb"
)

var initialQueueCapacity = 16

type queueRecord struct {
	cc.Queue
}

type queueRecordPtr = *queueRecord

func newQueueRecord() queueRecordPtr {
	return &queueRecord{cc.NewSynchronizedRingQueue(initialQueueCapacity)}
}

type queueMap struct {
	cc.Map
}

type queueMapPtr = *queueMap

func newQueueMap() queueMapPtr {
	return &queueMap{cc.NewSynchronizedMap(0)}
}

// MarshalBinary implements encoding.BinaryMarshaler.MarshalBinary
// (https://golang.org/pkg/encoding/#BinaryMarshaler).
func (m queueMapPtr) MarshalBinary() ([]byte, error) {
	// number of elements
	data := make([]byte, 4)
	sz := m.Size()
	binary.LittleEndian.PutUint32(data, uint32(sz))
	idx := 4

	m.Range(func(k, v interface{}) bool {
		n := k.(string)
		r := v.(queueRecordPtr)

		// name
		sz := len(n)
		data = append(data, make([]byte, 4+sz)...)
		binary.LittleEndian.PutUint32(data[idx:], uint32(sz))
		idx += 4
		idx += copy(data[idx:], []byte(n))

		// number of elements in the queue
		sz = r.Size()
		data = append(data, make([]byte, 4)...)
		binary.LittleEndian.PutUint32(data[idx:], uint32(sz))
		idx += 4

		// elements
		r.Range(func(v interface{}) bool {
			elm := v.(entryPtr)
			elmSz := len(elm.Value)

			// reserve space
			data = append(data, make([]byte, 4+elmSz+8)...)

			// value size
			binary.LittleEndian.PutUint32(data[idx:], uint32(elmSz))
			idx += 4

			// value
			idx += copy(data[idx:], elm.Value)

			// ttl
			binary.LittleEndian.PutUint64(data[idx:], uint64(elm.TTLMillis))
			idx += 8

			return true
		})

		return true
	})

	return data, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.UnmarshalBinary
// (https://golang.org/pkg/encoding/#BinaryUnmarshaler).
func (m queueMapPtr) UnmarshalBinary(data []byte) error {
	m.Clear()
	idx := 0

	// number of elements
	sz := int(binary.LittleEndian.Uint32(data[idx:]))
	if sz == 0 {
		return nil
	}
	idx += 4

	for i := 0; i < sz; i++ {
		// name
		sz := int(binary.LittleEndian.Uint32(data[idx:]))
		idx += 4
		b := make([]byte, sz)
		copy(b, data[idx:idx+sz])
		n := string(b)
		idx += sz

		// number of elements
		sz = int(binary.LittleEndian.Uint32(data[idx:]))
		idx += 4

		r := newQueueRecord()

		// elements
		for j := 0; j < sz; j++ {
			e := entry{}
			sz := int(binary.LittleEndian.Uint32(data[idx:]))
			idx += 4

			b := make([]byte, sz)
			copy(b, data[idx:idx+sz])
			e.Value = b
			idx += sz

			e.TTLMillis = int64(binary.LittleEndian.Uint64(data[idx:]))
			idx += 8

			r.Offer(&e)
		}

		m.Put(n, r)
	}

	return nil
}

func queue(f *FSM, n string) queueRecordPtr {
	v, _ := f.qs.ComputeIfAbsent(n, func() interface{} { return newQueueRecord() })
	return v.(queueRecordPtr)
}

func queueSize(f *FSM, m *pb.ClusterCommand) interface{} {
	q := queue(f, m.GetEntity().EntityName)
	v := make([]byte, 4)
	binary.LittleEndian.PutUint32(v, uint32(q.Size()))
	return &pb.Value{Value: v}
}

func queueClear(f *FSM, m *pb.ClusterCommand) interface{} {
	q := queue(f, m.GetEntity().EntityName)
	q.Clear()
	return nil
}

func queueOffer(f *FSM, m *pb.ClusterCommand) interface{} {
	v := m.GetValue()
	n := v.EntityName
	q := queue(f, n)
	q.Offer(&entry{v.Value, m.TtlMillis})
	return nil
}

func queuePoll(f *FSM, m *pb.ClusterCommand) interface{} {
	q := queue(f, m.GetEntity().EntityName)
	if v := q.Poll(); v != nil {
		return &pb.Value{Value: v.(entryPtr).Value}
	}
	return &pb.Value{}
}

func queuePeek(f *FSM, m *pb.ClusterCommand) interface{} {
	q := queue(f, m.GetEntity().EntityName)
	if v := q.Peek(); v != nil {
		return &pb.Value{Value: v.(entryPtr).Value}
	}
	return &pb.Value{}
}
