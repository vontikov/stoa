package cluster

import (
	"testing"

	"context"
	"sync"

	"github.com/stretchr/testify/assert"
)

func TestFSMQueue(t *testing.T) {
	const queueName = "test-queue"
	const max = 100

	f := NewFSM(context.Background())

	var wg sync.WaitGroup
	for i := 0; i < max; i++ {
		wg.Add(1)
		go func(n int) {
			q := queue(f, queueName)
			assert.NotNil(t, q)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestFSMQueueMarshalling(t *testing.T) {
	assert := assert.New(t)

	type qd struct {
		name    string
		entries []entry
	}

	tests := []struct {
		name             string
		queues           []qd
		expectedSize     int
		expectedDataSize int
	}{
		{
			name:             "No queues",
			expectedSize:     0,
			expectedDataSize: 4,
		},
		{
			name: "One empty queue",
			queues: []qd{
				{
					name: "abc",
				},
			},
			expectedSize: 1,
			expectedDataSize: 4 + // number of elements
				4 + 3 + // queue name
				4, // number of elements in the queue
		},
		{
			name: "One non-empty queue",
			queues: []qd{
				{
					name: "abc",
					entries: []entry{
						{[]byte{1, 2, 3}, 100},
						{[]byte{4, 5, 6}, 200},
					},
				},
			},
			expectedSize: 1,
			expectedDataSize: 4 + // number of elements
				4 + 3 + // queue name
				4 + // number of elements in the queue
				4 + 3 + 8 + // element 1
				4 + 3 + 8, // element 2
		},
		{
			name: "One empty and one non-empty queue",
			queues: []qd{
				{
					name: "abc",
				},
				{
					name: "xyz",
					entries: []entry{
						{[]byte{1, 2, 3}, 100},
						{[]byte{4, 5, 6}, 200},
					},
				},
			},
			expectedSize: 2,
			expectedDataSize: 4 + // number of elements
				4 + 3 + // queue name
				4 + // number of elements in the queue
				4 + 3 + // queue name
				4 + // number of elements in the queue
				4 + 3 + 8 + // element 1
				4 + 3 + 8, // element 2
		},

		{
			name: "Two non-empty queues",
			queues: []qd{
				{
					name: "abc",
					entries: []entry{
						{[]byte{1}, 100},
						{[]byte{2, 3}, 200},
						{[]byte{4, 5, 6}, 200},
					},
				},
				{
					name: "xyz",
					entries: []entry{
						{[]byte{1}, 100},
					},
				},
			},
			expectedSize: 2,
			expectedDataSize: 4 + // number of elements
				4 + 3 + // queue name
				4 + // number of elements in the queue
				4 + 1 + 8 + // element 1
				4 + 2 + 8 + // element 2
				4 + 3 + 8 + // element 3
				4 + 3 + // queue name
				4 + // number of elements in the queue
				4 + 1 + 8, // element 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newQueueMap()
			for _, qd := range tt.queues {
				q := newQueueRecord()
				for i := range qd.entries {
					e := &qd.entries[i]
					q.Offer(e)
				}
				m.Put(qd.name, q)
			}

			data, err := m.MarshalBinary()
			assert.Nil(err)
			assert.Equal(tt.expectedDataSize, len(data))

			dest := newQueueMap()
			err = dest.UnmarshalBinary(data)
			assert.Nil(err)
			assert.Equal(tt.expectedSize, dest.Size())

			for _, qd := range tt.queues {
				v := dest.Get(qd.name)
				assert.NotNil(v)
				q := v.(queueRecordPtr)
				assert.Equal(len(qd.entries), q.Size())
				for i := range qd.entries {
					assert.Equal(&qd.entries[i], q.Poll())
				}
			}
		})
	}
}
