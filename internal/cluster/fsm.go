package cluster

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"io"
	"sync"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/raft"

	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
)

var statusChanSize = 1024

// ErrUnknownCommand is the error returned by fsm.Apply if the command is
// unknown.
var ErrUnknownCommand = errors.New("unknown command")

// FSM implements raft.FSM
type FSM struct {
	logger logging.Logger
	ctx    context.Context

	mu sync.RWMutex // protects following fields
	qs queueMapPtr
	ds dictMapPtr
	ms mutexMapPtr

	statusMutex sync.Mutex // protects following fields
	statusChan  chan *pb.Status

	leadership int32 // must not be accessed directly
}

type fsmCommand func(*FSM, *pb.ClusterCommand) interface{}

// NewFSM creates a new raft.FSM instance
func NewFSM(ctx context.Context) *FSM {
	return &FSM{
		logger: logging.NewLogger("fsm"),
		ctx:    ctx,
		qs:     newQueueMap(),
		ds:     newDictMap(),
		ms:     newMutexMap(),
	}
}

var fsmCommands [pb.ClusterCommand_MAX_INDEX]fsmCommand = [...]fsmCommand{
	pb.ClusterCommand_QUEUE_SIZE:               queueSize,
	pb.ClusterCommand_QUEUE_CLEAR:              queueClear,
	pb.ClusterCommand_QUEUE_OFFER:              queueOffer,
	pb.ClusterCommand_QUEUE_POLL:               queuePoll,
	pb.ClusterCommand_QUEUE_PEEK:               queuePeek,
	pb.ClusterCommand_DICTIONARY_SIZE:          dictionarySize,
	pb.ClusterCommand_DICTIONARY_CLEAR:         dictionaryClear,
	pb.ClusterCommand_DICTIONARY_PUT:           dictionaryPut,
	pb.ClusterCommand_DICTIONARY_PUT_IF_ABSENT: dictionaryPutIfAbsent,
	pb.ClusterCommand_DICTIONARY_GET:           dictionaryGet,
	pb.ClusterCommand_DICTIONARY_REMOVE:        dictionaryRemove,
	pb.ClusterCommand_DICTIONARY_RANGE:         dictionaryRange,
	pb.ClusterCommand_MUTEX_TRY_LOCK:           mutexTryLock,
	pb.ClusterCommand_MUTEX_UNLOCK:             mutexUnlock,
	pb.ClusterCommand_SERVICE_PING:             processPing,
}

func (f *FSM) leader(v bool) {
	var t int32
	if v {
		t = 1
	}
	atomic.StoreInt32(&f.leadership, t)
}

func (f *FSM) isLeader() bool {
	return atomic.LoadInt32(&f.leadership) == 1
}

func (f *FSM) status() chan *pb.Status {
	if ch := f.statusChan; ch != nil {
		return ch
	}

	f.statusMutex.Lock()
	if f.statusChan == nil {
		f.statusChan = make(chan *pb.Status, statusChanSize)
	}
	f.statusMutex.Unlock()
	return f.statusChan
}

// Apply log is invoked once a log entry is committed.
func (f *FSM) Apply(l *raft.Log) interface{} {
	var c pb.ClusterCommand
	if err := proto.Unmarshal(l.Data, &c); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	return fsmCommands[c.Command](f, &c)
}

// Execute executes the command c.
func (f *FSM) Execute(c *pb.ClusterCommand) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	return fsmCommands[c.Command](f, c)
}

// Snapshot is used to support log compaction.
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return f, nil
}

// Persist implements raft.SnapshotSink.Persist.
func (f *FSM) Persist(sink raft.SnapshotSink) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	err := func() error {
		var b bytes.Buffer
		if err := gob.NewEncoder(&b).Encode(f); err != nil {
			return err
		}
		if _, err := sink.Write(b.Bytes()); err != nil {
			return err
		}
		return sink.Close()
	}()

	if err != nil {
		_ = sink.Cancel()
	}
	return err
}

// Restore is used to restore an FSM from a snapshot.
func (f *FSM) Restore(rc io.ReadCloser) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return gob.NewDecoder(rc).Decode(f)
}

// Release implements raft.SnapshotSink.Release.
func (f *FSM) Release() {}

// MarshalBinary implements encoding.BinaryMarshaler.MarshalBinary
// (https://golang.org/pkg/encoding/#BinaryMarshaler).
func (f *FSM) MarshalBinary() ([]byte, error) {
	f.logger.Debug("start marshalling")

	// queues
	qb, err := f.qs.MarshalBinary()
	if err != nil {
		return nil, err
	}
	qsz := len(qb)
	f.logger.Debug("queues marshalled", "size", qsz)

	// dictionaries
	db, err := f.ds.MarshalBinary()
	if err != nil {
		return nil, err
	}
	dsz := len(db)
	f.logger.Debug("dictionaries marshalled", "size", dsz)

	// mutexes
	mb, err := f.ms.MarshalBinary()
	if err != nil {
		return nil, err
	}
	msz := len(mb)
	f.logger.Debug("mutexes marshalled", "size", msz)

	// put together
	data := make([]byte, 4+qsz+4+dsz+4+msz)
	idx := 0

	// queues
	binary.LittleEndian.PutUint32(data[idx:], uint32(qsz))
	idx += 4
	idx += copy(data[idx:], qb)

	// dictionaries
	binary.LittleEndian.PutUint32(data[idx:], uint32(dsz))
	idx += 4
	idx += copy(data[idx:], db)

	// mutexes
	binary.LittleEndian.PutUint32(data[idx:], uint32(msz))
	idx += 4
	copy(data[idx:], mb)

	f.logger.Debug("marshalling complete", "size", len(data))
	return data, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.UnmarshalBinary
// (https://golang.org/pkg/encoding/#BinaryUnmarshaler).
func (f *FSM) UnmarshalBinary(data []byte) error {
	f.logger.Debug("start unmarshalling", "size", len(data))
	idx := 0

	// queues
	sz := int(binary.LittleEndian.Uint32(data[idx:]))
	idx += 4
	qs := newQueueMap()
	if err := qs.UnmarshalBinary(data[idx : idx+sz]); err != nil {
		return err
	}
	idx += sz
	f.logger.Debug("queues unmarshalled", "size", sz)

	// dictionaries
	sz = int(binary.LittleEndian.Uint32(data[idx:]))
	idx += 4
	ds := newDictMap()
	if err := ds.UnmarshalBinary(data[idx : idx+sz]); err != nil {
		return err
	}
	idx += sz
	f.logger.Debug("dictionaries unmarshalled", "size", sz)

	// mutexes
	sz = int(binary.LittleEndian.Uint32(data[idx:]))
	idx += 4
	ms := newMutexMap()
	if err := ms.UnmarshalBinary(data[idx : idx+sz]); err != nil {
		return err
	}
	f.logger.Debug("mutexes unmarshalled", "size", sz)

	f.qs = qs
	f.ds = ds
	f.ms = ms

	f.logger.Debug("unmarshalling complete")
	return nil
}
