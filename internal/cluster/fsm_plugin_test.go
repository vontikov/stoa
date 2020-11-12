package cluster

import (
	"testing"

	"context"
	"encoding/binary"
	"sort"
	"strconv"

	"github.com/stretchr/testify/assert"
	"github.com/vontikov/stoa/pkg/pb"
)

func TestFsmPlugin(t *testing.T) {
	const (
		max = 10
	)

	//	dicts := []string{"d1", "d2", "d3"}
	dicts := []string{"d1", "d2", "d3"}

	f := NewFSM(context.Background())

	// populate
	for _, n := range dicts {
		for i := 0; i < max; i++ {
			v := pb.KeyValue{
				Name:  n,
				Key:   []byte(strconv.Itoa(i)),
				Value: []byte(strconv.Itoa(i)),
			}
			m := pb.ClusterCommand{
				Payload: &pb.ClusterCommand_KeyValue{KeyValue: &v},
			}
			dictionaryPut(f, &m)
		}
		d := dictionary(f, n)
		assert.Equal(t, max, d.Size())
	}

	// scan
	r := make(map[string][]*pb.DictionaryEntry)
	for e := range f.DictionaryScan() {
		r[e.Name] = append(r[e.Name], e)
	}

	for _, n := range dicts {
		s := r[n]
		assert.Equal(t, max, len(s))

		var keys, vals []int
		for _, e := range s {
			k, _ := strconv.Atoi(string(e.Key))
			v, _ := strconv.Atoi(string(e.Value))
			keys = append(keys, k)
			vals = append(vals, v)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
		sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
		for i := 0; i < max; i++ {
			assert.Equal(t, i, keys[i])
			assert.Equal(t, i, vals[i])
		}
	}

	// remove
	for _, n := range dicts {
		m := &pb.ClusterCommand{
			Command: pb.ClusterCommand_DICTIONARY_SIZE,
			Payload: &pb.ClusterCommand_Name{Name: &pb.Name{Name: n}},
		}
		v := dictionarySize(f, m).(*pb.Value)
		assert.Equal(t, uint32(max), binary.LittleEndian.Uint32(v.Value))
	}
	for e := range f.DictionaryScan() {
		f.DictionaryRemove(e)
	}
	for _, n := range dicts {
		m := &pb.ClusterCommand{
			Command: pb.ClusterCommand_DICTIONARY_SIZE,
			Payload: &pb.ClusterCommand_Name{Name: &pb.Name{Name: n}},
		}
		v := dictionarySize(f, m).(*pb.Value)
		assert.Equal(t, uint32(0), binary.LittleEndian.Uint32(v.Value))
	}
}
