package cluster

import (
	"testing"

	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vontikov/stoa/pkg/pb"
)

func TestFsmMutex(t *testing.T) {
	const muxName = "test-mutex"
	const max = 100

	f := NewFSM(context.Background())

	var wg sync.WaitGroup
	for i := 0; i < max; i++ {
		wg.Add(1)
		go func(n int) {
			mx := mutex(f, muxName)
			assert.NotNil(t, mx)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestFsmMutexTryLockUnlock(t *testing.T) {
	assert := assert.New(t)

	id := []byte("test-id")
	const muxName = "test-mutex"
	const max = 100

	f := NewFSM(context.Background())

	m := &pb.ClusterCommand{
		Command: pb.ClusterCommand_MUTEX_TRY_LOCK,
		Payload: &pb.ClusterCommand_ClientId{ClientId: &pb.ClientId{Id: id}},
	}

	r := mutexTryLock(f, m).(*pb.Result)
	assert.True(r.Ok)

	var wg sync.WaitGroup
	for i := 0; i < max; i++ {
		wg.Add(1)
		go func(n int) {
			m := &pb.ClusterCommand{
				Command: pb.ClusterCommand_MUTEX_TRY_LOCK,
				Payload: &pb.ClusterCommand_ClientId{ClientId: &pb.ClientId{Id: []byte(uuid.New().String())}},
			}

			r := mutexTryLock(f, m).(*pb.Result)
			assert.False(r.Ok)

			r = mutexUnlock(f, m).(*pb.Result)
			assert.False(r.Ok)

			wg.Done()
		}(i)
	}
	wg.Wait()

	r = mutexUnlock(f, m).(*pb.Result)
	assert.True(r.Ok)
}

func TestFSMMutexMarshalling(t *testing.T) {
	assert := assert.New(t)

	type md struct {
		name string
		r    mutexRecord
	}

	tests := []struct {
		name             string
		mutexes          []md
		expectedDataSize int
		expectedSize     int
	}{
		{
			name:             "No mutexes",
			mutexes:          nil,
			expectedDataSize: 4,
			expectedSize:     0,
		},
		{
			name: "One unlocked mutex",
			mutexes: []md{
				{"mx-name", mutexRecord{}},
			},
			expectedDataSize: 4 + // number of elements
				4 + 7 + // mutex name
				2 + // mutex locked
				4 + // mutex lockedBy (empty)
				4 + 15, // touched
			expectedSize: 1,
		},
		{
			name: "One locked mutex",
			mutexes: []md{
				{"mx-name", mutexRecord{locked: true, lockedBy: []byte("client")}},
			},
			expectedDataSize: 4 + // number of elements
				4 + 7 + // mutex name
				2 + // mutex locked
				4 + 6 + // mutex lockedBy (non-empty)
				4 + 15, // touched
			expectedSize: 1,
		},
		{
			name: "One locked and touched",
			mutexes: []md{
				{"mx-name", mutexRecord{locked: true, lockedBy: []byte("client"), touched: time.Now()}},
			},
			expectedDataSize: 4 + // number of elements
				4 + 7 + // mutex name
				2 + // mutex locked
				4 + 6 + // mutex lockedBy (non-empty)
				4 + 15, // touched
			expectedSize: 1,
		},
		{
			name: "One unlocked and one locked and touched",
			mutexes: []md{
				{"mxname0", mutexRecord{}},
				{"mxname1", mutexRecord{locked: true, lockedBy: []byte("client"), touched: time.Now()}},
			},
			expectedDataSize: 4 + // number of elements
				// first
				4 + 7 + // mutex name
				2 + // mutex locked
				4 + // mutex lockedBy (empty)
				4 + 15 + // touched
				// second
				4 + 7 + // mutex name
				2 + // mutex locked
				4 + 6 + // mutex lockedBy (non-empty)
				4 + 15, // touched
			expectedSize: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMutexMap()
			for i := range tt.mutexes {
				md := &tt.mutexes[i]
				m.Put(md.name, &md.r)
			}

			data, err := m.MarshalBinary()
			assert.Nil(err)
			assert.Equal(tt.expectedDataSize, len(data))

			dest := newMutexMap()
			err = dest.UnmarshalBinary(data)
			assert.Nil(err)
			assert.Equal(tt.expectedSize, dest.Size())

			for i := range tt.mutexes {
				md := &tt.mutexes[i]
				v := dest.Get(md.name)
				assert.NotNil(v)
				r := v.(mutexRecordPtr)
				assert.Equal(md.r.locked, r.locked)
				assert.Equal(md.r.lockedBy, r.lockedBy)
			}
		})
	}
}

func TestMutexExpiration(t *testing.T) {
	assert := assert.New(t)

	const (
		muxName = "test-mutex"
	)

	clientId := []byte("abc")

	f := NewFSM(context.Background())
	f.leader(true)

	mx := mutex(f, muxName)
	assert.NotNil(mx)

	r := mx.tryLock(clientId)
	assert.True(r)

	expiration := 300 * time.Millisecond

	time.Sleep(expiration / 2)
	n := mutexUnlockExpired(f, expiration)
	assert.Equal(0, n, "should not unlock yet")
	assert.True(mx.isLocked())
	assert.Equal(0, len(f.status()))

	time.Sleep(expiration)
	n = mutexUnlockExpired(f, expiration)
	assert.Equal(1, n, "should unlock")
	assert.False(mx.isLocked())
	assert.Equal(1, len(f.status()))

	s := <-f.status()
	assert.Equal(muxName, s.GetM().Name)
	assert.False(s.GetM().Locked)
}
