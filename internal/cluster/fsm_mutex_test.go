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
			mx := f.mutex(muxName)
			assert.NotNil(t, mx)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestFsmMutexTryLockUnlock(t *testing.T) {
	const muxName = "test-mutex"
	const id = "test-id"
	const max = 100

	f := NewFSM(context.Background())

	m := &pb.ClusterCommand{
		Command: pb.ClusterCommand_MUTEX_TRY_LOCK,
		Payload: &pb.ClusterCommand_Id{Id: &pb.Id{Id: id}},
	}

	r := f.mutexTryLock(m).(*pb.Result)
	assert.True(t, r.Ok)

	var wg sync.WaitGroup
	for i := 0; i < max; i++ {
		wg.Add(1)
		go func(n int) {
			m := &pb.ClusterCommand{
				Command: pb.ClusterCommand_MUTEX_TRY_LOCK,
				Payload: &pb.ClusterCommand_Id{Id: &pb.Id{Id: uuid.New().String()}},
			}

			r := f.mutexTryLock(m).(*pb.Result)
			assert.False(t, r.Ok)

			r = f.mutexUnlock(m).(*pb.Result)
			assert.False(t, r.Ok)

			wg.Done()
		}(i)
	}
	wg.Wait()

	r = f.mutexUnlock(m).(*pb.Result)
	assert.True(t, r.Ok)
}

func TestFsmMutexWatcher(t *testing.T) {
	t.Run("Should expire", func(t *testing.T) {
		mutexCheckPeriod = 10 * time.Millisecond
		mutexDeadline = 50 * time.Millisecond

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		f := NewFSM(ctx)
		mx := f.mutex("muxName")
		r := mx.tryLock("lockId")
		assert.True(t, r)
		time.Sleep(mutexDeadline << 1)
		assert.False(t, mx.isLocked())
	})

	t.Run("Should not expire", func(t *testing.T) {
		mutexCheckPeriod = 50 * time.Millisecond
		mutexDeadline = 10 * time.Millisecond

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		f := NewFSM(ctx)
		mx := f.mutex("muxName")
		r := mx.tryLock("lockId")
		assert.True(t, r)
		time.Sleep(mutexDeadline << 1)
		assert.True(t, mx.isLocked())
	})
}
