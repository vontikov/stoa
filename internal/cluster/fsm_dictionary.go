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

func (f *FSM) dictionary(n string) cc.Map {
	v, _ := f.ds.ComputeIfAbsent(n, func() interface{} { return newDictionaryRecord() })
	return v.(*dictionaryRecord).m
}

func (f *FSM) dictionarySize(m *pb.ClusterCommand) interface{} {
	d := f.dictionary(m.GetName().Name)
	v := make([]byte, 4)
	binary.LittleEndian.PutUint32(v, uint32(d.Size()))
	return &pb.Value{Value: v}
}

func (f *FSM) dictionaryClear(m *pb.ClusterCommand) interface{} {
	d := f.dictionary(m.GetName().Name)
	d.Clear()
	return nil
}

func (f *FSM) dictionaryPut(m *pb.ClusterCommand) interface{} {
	kv := m.GetKeyValue()
	d := f.dictionary(kv.Name)
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

func (f *FSM) dictionaryPutIfAbsent(m *pb.ClusterCommand) interface{} {
	kv := m.GetKeyValue()
	d := f.dictionary(kv.Name)
	v := entry{
		Value:     kv.Value,
		TTLMillis: m.TtlMillis,
	}
	ok := d.PutIfAbsent(string(kv.Key), v)
	return &pb.Result{Ok: ok}
}

func (f *FSM) dictionaryGet(m *pb.ClusterCommand) interface{} {
	k := m.GetKey()
	d := f.dictionary(k.Name)
	if v := d.Get(string(k.Key)); v != nil {
		return &pb.Value{Value: v.(entry).Value}
	}
	return &pb.Value{}
}

func (f *FSM) dictionaryRemove(m *pb.ClusterCommand) interface{} {
	k := m.GetKey()
	d := f.dictionary(k.Name)
	d.Remove(string(k.Key))
	return nil
}

func (f *FSM) dictionaryScan(m *pb.ClusterCommand) interface{} {
	ch := make(chan *pb.KeyValue)

	go func() {
		d := f.dictionary(m.GetName().Name)
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
