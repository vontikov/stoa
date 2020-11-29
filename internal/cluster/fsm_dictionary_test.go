package cluster

import (
	"testing"

	"context"
	"sync"

	"github.com/stretchr/testify/assert"
)

func TestFsmDictionary(t *testing.T) {
	const dictName = "test-dictionary"
	const max = 100

	f := NewFSM(context.Background())

	var wg sync.WaitGroup
	for i := 0; i < max; i++ {
		wg.Add(1)
		go func(n int) {
			d := dictionary(f, dictName)
			assert.NotNil(t, d)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestFSMDictionaryMarshalling(t *testing.T) {
	assert := assert.New(t)

	type kv struct {
		key []byte
		val entry
	}

	type dd struct {
		name    string
		entries []kv
	}

	tests := []struct {
		name             string
		dicts            []dd
		expectedSize     int
		expectedDataSize int
	}{
		{
			name:             "No dicts",
			expectedSize:     0,
			expectedDataSize: 4,
		},
		{
			name: "One empty dictionary",
			dicts: []dd{
				{
					name: "abc",
				},
			},
			expectedSize: 1,
			expectedDataSize: 4 + // number of elements
				4 + 3 + // dictionary name
				4, // number of elements in the dictionary
		},
		{
			name: "One non-empty dictionary",
			dicts: []dd{
				{
					name: "abc",
					entries: []kv{
						{[]byte{1, 2, 3}, entry{[]byte{1, 2, 3}, 100}},
						{[]byte{4, 5, 6, 7}, entry{[]byte{4, 5, 6}, 200}},
					},
				},
			},
			expectedSize: 1,
			expectedDataSize: 4 + // number of elements
				4 + 3 + // dictionary name
				4 + // number of elements in the dictionary
				4 + 3 + 4 + 3 + 8 + // element 1
				4 + 4 + 4 + 3 + 8, // element 2
		},
		{
			name: "One empty and one non-empty dictionary",
			dicts: []dd{
				{
					name: "abc",
				},
				{
					name: "xyz",
					entries: []kv{
						{[]byte{1, 2, 3}, entry{[]byte{1, 2, 3}, 100}},
						{[]byte{4, 5, 6, 7}, entry{[]byte{4, 5, 6}, 200}},
					},
				},
			},
			expectedSize: 2,
			expectedDataSize: 4 + // number of elements
				4 + 3 + // dictionary name
				4 + // number of elements in the dictionary
				4 + 3 + // dictionary name
				4 + // number of elements in the dictionary
				4 + 3 + 4 + 3 + 8 + // element 1
				4 + 4 + 4 + 3 + 8, // element 2
		},

		{
			name: "Two non-empty dictionary",
			dicts: []dd{
				{
					name: "abc",
					entries: []kv{
						{[]byte{1, 2, 3}, entry{[]byte{1, 2, 3}, 100}},
						{[]byte{4, 5, 6, 7}, entry{[]byte{4, 5, 6}, 200}},
					},
				},
				{
					name: "xyz",
					entries: []kv{
						{[]byte{1, 2, 3}, entry{[]byte{1, 2, 3}, 100}},
						{[]byte{4, 5}, entry{[]byte{4, 5, 6}, 200}},
					},
				},
			},
			expectedSize: 2,
			expectedDataSize: 4 + // number of elements
				4 + 3 + // dictionary name
				4 + // number of elements in the dictionary
				4 + 3 + 4 + 3 + 8 + // element 1
				4 + 4 + 4 + 3 + 8 + // element 2
				4 + 3 + // dictionary name
				4 + // number of elements in the dictionary
				4 + 3 + 4 + 3 + 8 + // element 1
				4 + 2 + 4 + 3 + 8, // element 2
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newDictMap()
			for _, dd := range tt.dicts {
				d := newDictRecord()
				for i := range dd.entries {
					e := dd.entries[i]
					d.Put(string(e.key), &e.val)
				}
				m.Put(dd.name, d)
			}

			data, err := m.MarshalBinary()
			assert.Nil(err)
			assert.Equal(tt.expectedDataSize, len(data))

			dest := newDictMap()
			err = dest.UnmarshalBinary(data)
			assert.Nil(err)
			assert.Equal(tt.expectedSize, dest.Size())

			for _, dd := range tt.dicts {
				v := dest.Get(dd.name)
				assert.NotNil(v)
				d := v.(dictRecordPtr)
				assert.Equal(len(dd.entries), d.Size())
				for i := range dd.entries {
					e := dd.entries[i]
					assert.Equal(&e.val, d.Get(string(e.key)))
				}
			}
		})
	}
}
