package cluster

import (
	"encoding/binary"

	cc "github.com/vontikov/go-concurrent"
	"github.com/vontikov/stoa/pkg/pb"
)

var initialDictionaryCapacity = 16

type dictRecord struct {
	cc.Map
}

type dictRecordPtr = *dictRecord

func newDictRecord() dictRecordPtr {
	return &dictRecord{cc.NewSynchronizedMap(initialDictionaryCapacity)}
}

type dictMap struct {
	cc.Map
}

type dictMapPtr = *dictMap

func newDictMap() dictMapPtr {
	return &dictMap{cc.NewSynchronizedMap(0)}
}

// MarshalBinary implements encoding.BinaryMarshaler.MarshalBinary
// (https://golang.org/pkg/encoding/#BinaryMarshaler).
func (m dictMapPtr) MarshalBinary() ([]byte, error) {
	data := make([]byte, 4)
	sz := m.Size()
	// number of elements
	binary.LittleEndian.PutUint32(data, uint32(sz))
	idx := 4

	m.Range(func(k, v interface{}) bool {
		n := k.(string)
		d := v.(dictRecordPtr)

		// dictionary name
		sz := len(n)
		data = append(data, make([]byte, 4+sz)...)
		binary.LittleEndian.PutUint32(data[idx:], uint32(sz))
		idx += 4
		idx += copy(data[idx:], []byte(n))

		// number of elements in the dictionary
		sz = d.Size()
		data = append(data, make([]byte, 4)...)
		binary.LittleEndian.PutUint32(data[idx:], uint32(sz))
		idx += 4

		// elements
		d.Range(func(k, v interface{}) bool {
			key := []byte(k.(string))
			keySz := len(key)

			val := v.(entryPtr)
			valSz := len(val.Value)

			// reserve space
			data = append(data, make([]byte, 4+keySz+4+valSz+8)...)

			// key size
			binary.LittleEndian.PutUint32(data[idx:], uint32(keySz))
			idx += 4

			// key
			idx += copy(data[idx:], key)

			// value size
			binary.LittleEndian.PutUint32(data[idx:], uint32(valSz))
			idx += 4

			// value
			idx += copy(data[idx:], val.Value)

			// ttl
			binary.LittleEndian.PutUint64(data[idx:], uint64(val.TTLMillis))
			idx += 8

			return true
		})
		return true
	})

	return data, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.UnmarshalBinary
// (https://golang.org/pkg/encoding/#BinaryUnmarshaler).
func (m dictMapPtr) UnmarshalBinary(data []byte) error {
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

		d := newDictRecord()

		// elements
		for j := 0; j < sz; j++ {
			// key
			sz := int(binary.LittleEndian.Uint32(data[idx:]))
			idx += 4
			b := make([]byte, sz)
			copy(b, data[idx:idx+sz])
			key := b
			idx += sz

			// value
			val := entry{}
			sz = int(binary.LittleEndian.Uint32(data[idx:]))
			idx += 4

			b = make([]byte, sz)
			copy(b, data[idx:idx+sz])
			val.Value = b
			idx += sz

			val.TTLMillis = int64(binary.LittleEndian.Uint64(data[idx:]))
			idx += 8

			d.Put(string(key), &val)
		}

		m.Put(n, d)
	}

	return nil
}

func dictionary(f *FSM, n string) dictRecordPtr {
	v, _ := f.ds.ComputeIfAbsent(n, func() interface{} { return newDictRecord() })
	return v.(dictRecordPtr)
}

func dictionarySize(f *FSM, m *pb.ClusterCommand) interface{} {
	d := dictionary(f, m.GetEntity().EntityName)
	v := make([]byte, 4)
	binary.LittleEndian.PutUint32(v, uint32(d.Size()))
	return &pb.Value{Value: v}
}

func dictionaryClear(f *FSM, m *pb.ClusterCommand) interface{} {
	d := dictionary(f, m.GetEntity().EntityName)
	d.Clear()
	return nil
}

func dictionaryPut(f *FSM, m *pb.ClusterCommand) interface{} {
	kv := m.GetKeyValue()
	d := dictionary(f, kv.EntityName)
	v := entry{
		Value:     kv.Value,
		TTLMillis: m.TtlMillis,
	}
	o := d.Put(string(kv.Key), &v)
	if o == nil {
		return &pb.Value{}
	}
	return &pb.Value{Value: o.(*entry).Value}
}

func dictionaryPutIfAbsent(f *FSM, m *pb.ClusterCommand) interface{} {
	kv := m.GetKeyValue()
	d := dictionary(f, kv.EntityName)
	v := entry{
		Value:     kv.Value,
		TTLMillis: m.TtlMillis,
	}
	ok := d.PutIfAbsent(string(kv.Key), &v)
	return &pb.Result{Ok: ok}
}

func dictionaryGet(f *FSM, m *pb.ClusterCommand) interface{} {
	k := m.GetKey()
	d := dictionary(f, k.EntityName)
	if v := d.Get(string(k.Key)); v != nil {
		return &pb.Value{Value: v.(*entry).Value}
	}
	return &pb.Value{}
}

func dictionaryRemove(f *FSM, m *pb.ClusterCommand) interface{} {
	k := m.GetKey()
	d := dictionary(f, k.EntityName)
	d.Remove(string(k.Key))
	return nil
}

func dictionaryRange(f *FSM, m *pb.ClusterCommand) interface{} {
	ch := make(chan *pb.KeyValue)
	go func() {
		n := m.GetEntity().EntityName
		if f.logger.IsTrace() {
			f.logger.Trace("dictionary scan", "name", n)
		}

		d := dictionary(f, n)
		b := make([]*pb.KeyValue, 0, d.Size())
		d.Range(func(k, v interface{}) bool {
			kv := pb.KeyValue{
				Key:   []byte(k.(string)),
				Value: v.(*entry).Value,
			}
			b = append(b, &kv)
			return true
		})
		for _, kv := range b {
			ch <- kv
		}
		close(ch)
	}()
	return ch
}
