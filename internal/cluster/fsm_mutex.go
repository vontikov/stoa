package cluster

import (
	"sync"
	"time"

	"github.com/vontikov/stoa/pkg/pb"
)

var timeNow = time.Now
var mutexCheckPeriod = 200 * time.Millisecond
var mutexDeadline = 5 * time.Second

type mutexRecord struct {
	sync.RWMutex
	locked   bool
	lockedBy string
	touched  time.Time
}

func newMutexRecord() *mutexRecord {
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

func mutex(f *FSM, n string) *mutexRecord {
	v, ok := f.ms.ComputeIfAbsent(n, func() interface{} { return newMutexRecord() })
	if ok {
		f.mo.Do(f.startMutexWatcher)
	}
	return v.(*mutexRecord)
}

func mutexTryLock(f *FSM, m *pb.ClusterCommand) interface{} {
	id := m.GetId()
	mx := mutex(f, id.Name)
	r := mx.tryLock(id.Id)
	if r {
		notifyMutex(f, id.Name, mx)
	}
	return &pb.Result{Ok: r}
}

func mutexUnlock(f *FSM, m *pb.ClusterCommand) interface{} {
	id := m.GetId()
	mx := mutex(f, id.Name)
	r := mx.unlock(id.Id)
	if r {
		notifyMutex(f, id.Name, mx)
	}
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
		mx := v.(*mutexRecord)
		if mx.locked && mx.touched.Before(deadline) {
			f.logger.Warn("Mutex expired",
				"locked", mx.lockedBy, "locked by", mx.lockedBy, "touched", mx.touched)
			mx.unlock(mx.lockedBy)
			notifyMutex(f, k.(string), mx)
		}
	}
}

func notifyMutex(f *FSM, name string, mx *mutexRecord) {
	keys := f.streams.Keys()
	for _, k := range keys {
		v := f.streams.Get(k)
		stream := v.(pb.Stoa_KeepServer)

		status := pb.Status{Payload: &pb.Status_Mutex{Mutex: &pb.MutexStatus{Name: name, Locked: mx.locked}}}
		err := stream.Send(&status)
		if err != nil {
			f.logger.Error("Notification error", "message", err)
		}
	}
}
