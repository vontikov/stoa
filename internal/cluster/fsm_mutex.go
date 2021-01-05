package cluster

import (
	"sync"
	"time"

	cc "github.com/vontikov/go-concurrent"

	"github.com/vontikov/stoa/pkg/pb"
)

var timeNow = time.Now
var mutexCheckPeriod = 200 * time.Millisecond
var mutexDeadline = 5 * time.Second

type mutexMap struct {
	cc.Map
	once sync.Once
}

type mutexMapPtr = *mutexMap

func newMutexMap() mutexMapPtr {
	return &mutexMap{Map: cc.NewSynchronizedMap(0)}
}

type mutexRecord struct {
	sync.RWMutex
	locked   bool
	lockedBy string
	touched  time.Time
}

type mutexRecordPtr = *mutexRecord

func newMutexRecord() mutexRecordPtr {
	return &mutexRecord{}
}

func (m *mutexRecord) tryLock(id string) bool {
	m.Lock()
	defer m.Unlock()
	if m.locked {
		return false
	}
	m.locked = true
	m.lockedBy = id
	m.touched = timeNow()
	return true
}

func (m *mutexRecord) unlock(id string) bool {
	m.Lock()
	defer m.Unlock()
	if !m.locked || m.lockedBy != id {
		return false
	}
	m.locked = false
	m.lockedBy = ""
	return true
}

func (m *mutexRecord) isLocked() bool {
	m.Lock()
	defer m.Unlock()
	return m.locked
}

func (m *mutexRecord) touch(id string) bool {
	m.Lock()
	defer m.Unlock()
	if !m.locked || m.lockedBy != id {
		return false
	}
	m.touched = timeNow()
	return true
}

func (f *FSM) startMutexWatcher() {
	go func() {
		t := time.NewTicker(mutexCheckPeriod)
		for {
			select {
			case <-f.ctx.Done():
				return
			case <-t.C:
				checkExpiredMutexes(f)
			}
		}
	}()
	f.logger.Debug("Mutex watcher started")
}

func mutex(f *FSM, n string) mutexRecordPtr {
	v, ok := f.ms.ComputeIfAbsent(n, func() interface{} { return newMutexRecord() })
	if ok {
		f.ms.once.Do(f.startMutexWatcher)
	}
	return v.(mutexRecordPtr)
}

func mutexTryLock(f *FSM, m *pb.ClusterCommand) interface{} {
	id := m.GetId()
	mx := mutex(f, id.Name)
	r := mx.tryLock(id.Id)
	return &pb.Result{Ok: r}
}

func mutexUnlock(f *FSM, m *pb.ClusterCommand) interface{} {
	id := m.GetId()
	mx := mutex(f, id.Name)
	r := mx.unlock(id.Id)
	return &pb.Result{Ok: r}
}

func checkExpiredMutexes(f *FSM) {
	deadline := timeNow().Add(-mutexDeadline)
	keys := f.ms.Keys()
	for _, k := range keys {
		v := f.ms.Get(k)
		if v == nil {
			continue
		}
		mx := v.(mutexRecordPtr)
		if mx.locked && mx.touched.Before(deadline) {
			f.logger.Warn("Mutex expired",
				"locked", mx.lockedBy, "locked by", mx.lockedBy, "touched", mx.touched)
			mx.unlock(mx.lockedBy)
			// TODO notifyMutex()
		}
	}
}
