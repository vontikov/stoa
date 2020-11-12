package cluster

import (
	"github.com/vontikov/stoa/pkg/pb"
)

func (f *FSM) DictionaryScan() chan *pb.DictionaryEntry {
	b := make([]*pb.DictionaryEntry, 0, 1024)
	f.ds.Range(func(k, v interface{}) bool {
		n := k.(string)
		d := v.(*dictionaryRecord).m
		d.Range(func(k, v interface{}) bool {
			e := pb.DictionaryEntry{
				Name:      n,
				Key:       []byte(k.(string)),
				Value:     v.(entry).Value,
				TtlMillis: v.(entry).TTLMillis,
			}
			b = append(b, &e)
			return true
		})
		return true
	})

	ch := make(chan *pb.DictionaryEntry)
	go func() {
		for _, e := range b {
			ch <- e
		}
		close(ch)
	}()
	return ch
}

func (f *FSM) DictionaryRemove(e *pb.DictionaryEntry) error {
	d := dictionary(f, e.Name)
	d.Remove(string(e.Key))
	return nil
}
