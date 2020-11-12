package cluster

import (
	"encoding/binary"

	cc "github.com/vontikov/go-concurrent"
	"github.com/vontikov/stoa/pkg/pb"
)

var initialDictionaryCapacity = 16

type dictionaryRecord struct {
	m cc.Map
}

func newDictionaryRecord() *dictionaryRecord {
	return &dictionaryRecord{m: cc.NewSynchronizedMap(initialDictionaryCapacity)}
}

func dictionary(f *FSM, n string) cc.Map {
	v, _ := f.ds.ComputeIfAbsent(n, func() interface{} { return newDictionaryRecord() })
	return v.(*dictionaryRecord).m
}

func dictionarySize(f *FSM, m *pb.ClusterCommand) interface{} {
	d := dictionary(f, m.GetName().Name)
	v := make([]byte, 4)
	binary.LittleEndian.PutUint32(v, uint32(d.Size()))
	return &pb.Value{Value: v}
}

func dictionaryClear(f *FSM, m *pb.ClusterCommand) interface{} {
	d := dictionary(f, m.GetName().Name)
	d.Clear()
	return nil
}

func dictionaryPut(f *FSM, m *pb.ClusterCommand) interface{} {
	kv := m.GetKeyValue()
	d := dictionary(f, kv.Name)
	v := entry{
		Value:     kv.Value,
		TTLMillis: m.TtlMillis,
	}
	o := d.Put(string(kv.Key), v)
	if o == nil {
		return &pb.Value{}
	}
	return &pb.Value{Value: o.(entry).Value}
}

func dictionaryPutIfAbsent(f *FSM, m *pb.ClusterCommand) interface{} {
	kv := m.GetKeyValue()
	d := dictionary(f, kv.Name)
	v := entry{
		Value:     kv.Value,
		TTLMillis: m.TtlMillis,
	}
	ok := d.PutIfAbsent(string(kv.Key), v)
	return &pb.Result{Ok: ok}
}

func dictionaryGet(f *FSM, m *pb.ClusterCommand) interface{} {
	k := m.GetKey()
	d := dictionary(f, k.Name)
	if v := d.Get(string(k.Key)); v != nil {
		return &pb.Value{Value: v.(entry).Value}
	}
	return &pb.Value{}
}

func dictionaryRemove(f *FSM, m *pb.ClusterCommand) interface{} {
	k := m.GetKey()
	d := dictionary(f, k.Name)
	d.Remove(string(k.Key))
	return nil
}

func dictionaryScan(f *FSM, m *pb.ClusterCommand) interface{} {
	ch := make(chan *pb.KeyValue)
	go func() {
		d := dictionary(f, m.GetName().Name)
		b := make([]*pb.KeyValue, 0, d.Size())
		d.Range(func(k, v interface{}) bool {
			kv := pb.KeyValue{
				Key:   []byte(k.(string)),
				Value: v.(entry).Value,
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
