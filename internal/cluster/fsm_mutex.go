package cluster

import (
	"bytes"
	"encoding/binary"
	"sync"
	"time"

	cc "github.com/vontikov/go-concurrent"

	"github.com/vontikov/stoa/pkg/pb"
)

type mutexMap struct {
	cc.Map
}

type mutexMapPtr = *mutexMap

func newMutexMap() mutexMapPtr {
	return &mutexMap{Map: cc.NewSynchronizedMap(0)}
}

// MarshalBinary implements encoding.BinaryMarshaler.MarshalBinary
// (https://golang.org/pkg/encoding/#BinaryMarshaler).
func (m mutexMapPtr) MarshalBinary() ([]byte, error) {
	// number of elements
	data := make([]byte, 4)
	sz := m.Size()
	binary.LittleEndian.PutUint32(data, uint32(sz))
	idx := 4

	var err error
	m.Range(func(k, v interface{}) bool {
		n := k.(string)
		r := v.(mutexRecordPtr)

		// name
		sz := len(n)
		data = append(data, make([]byte, 4+sz)...)
		binary.LittleEndian.PutUint32(data[idx:], uint32(sz))
		idx += 4
		idx += copy(data[idx:], []byte(n))

		// locked
		var locked uint16
		if r.locked {
			locked = 1
		}
		data = append(data, make([]byte, 2)...)
		binary.LittleEndian.PutUint16(data[idx:], locked)
		idx += 2

		// lockedBy
		sz = len(r.lockedBy)
		data = append(data, make([]byte, 4+sz)...)
		binary.LittleEndian.PutUint32(data[idx:], uint32(sz))
		idx += 4
		idx += copy(data[idx:], []byte(r.lockedBy))

		// touched
		var touched []byte
		touched, err = r.touched.MarshalBinary()
		if err != nil {
			return false
		}
		sz = len(touched)
		data = append(data, make([]byte, 4+sz)...)
		binary.LittleEndian.PutUint32(data[idx:], uint32(sz))
		idx += 4
		idx += copy(data[idx:], touched)

		return true
	})

	if err != nil {
		return nil, err
	}

	return data, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.UnmarshalBinary
// (https://golang.org/pkg/encoding/#BinaryUnmarshaler).
func (m mutexMapPtr) UnmarshalBinary(data []byte) error {
	m.Clear()
	idx := 0

	// number of elements
	sz := int(binary.LittleEndian.Uint32(data[idx:]))
	if sz == 0 {
		return nil
	}
	idx += 4

	for i := 0; i < sz; i++ {
		// name
		sz := int(binary.LittleEndian.Uint32(data[idx:]))
		idx += 4
		b := make([]byte, sz)
		copy(b, data[idx:idx+sz])
		n := string(b)
		idx += sz

		r := newMutexRecord()

		// locked
		locked := binary.LittleEndian.Uint16(data[idx:])
		idx += 2
		if locked == 1 {
			r.locked = true
		}

		// lockedBy
		sz = int(binary.LittleEndian.Uint32(data[idx:]))
		idx += 4
		if sz > 0 {
			b = make([]byte, sz)
			copy(b, data[idx:idx+sz])
			r.lockedBy = b
			idx += sz
		}

		// touched
		sz = int(binary.LittleEndian.Uint32(data[idx:]))
		idx += 4
		if err := r.touched.UnmarshalBinary(data[idx : idx+sz]); err != nil {
			return err
		}
		idx += sz

		m.Put(n, r)
	}

	return nil
}

type mutexRecord struct {
	sync.RWMutex
	locked   bool
	lockedBy []byte
	touched  time.Time
	payload  []byte
}

type mutexRecordPtr = *mutexRecord

func newMutexRecord() mutexRecordPtr {
	return &mutexRecord{}
}

func (m *mutexRecord) tryLock(id []byte, payload []byte) (bool, []byte) {
	m.Lock()
	defer m.Unlock()
	if m.locked {
		return false, m.payload
	}
	m.locked = true
	m.lockedBy = id
	m.touched = timeNow()
	m.payload = payload
	return true, nil
}

func (m *mutexRecord) unlock(id []byte) (bool, []byte) {
	m.Lock()
	defer m.Unlock()
	if !m.locked || !bytes.Equal(m.lockedBy, id) {
		return false, nil
	}
	p := m.payload
	m.locked = false
	m.lockedBy = nil
	m.payload = nil
	return true, p
}

func (m *mutexRecord) isLocked() bool {
	m.Lock()
	defer m.Unlock()
	return m.locked
}

func (m *mutexRecord) touch(id []byte) bool {
	m.Lock()
	defer m.Unlock()
	if !m.locked || !bytes.Equal(m.lockedBy, id) {
		return false
	}
	m.touched = timeNow()
	return true
}

func mutex(f *FSM, n string) mutexRecordPtr {
	v, _ := f.ms.ComputeIfAbsent(n, func() interface{} { return newMutexRecord() })
	return v.(mutexRecordPtr)
}

func mutexTryLock(f *FSM, m *pb.ClusterCommand) interface{} {
	clientId := m.GetClientId()

	n := clientId.EntityName
	mx := mutex(f, n)

	id := clientId.Id
	r, p := mx.tryLock(id, clientId.Payload)
	return &pb.Result{Ok: r, Payload: p}
}

func mutexUnlock(f *FSM, m *pb.ClusterCommand) interface{} {
	clientId := m.GetClientId()

	n := clientId.EntityName
	mx := mutex(f, n)

	id := clientId.Id
	r, p := mx.unlock(id)
	return &pb.Result{Ok: r, Payload: p}
}

func mutexUnlockExpired(f *FSM, expiration time.Duration) int {
	n := 0
	deadline := timeNow().Add(-expiration)
	f.ms.Range(func(k, v interface{}) bool {
		m := v.(mutexRecordPtr)
		m.Lock()
		defer m.Unlock()

		if m.locked && m.touched.Before(deadline) {
			m.locked = false
			m.lockedBy = nil
			n++
			f.logger.Debug("expired mutex unlocked", "name", k)
		}
		return true
	})
	return n
}
