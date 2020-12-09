package cluster

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"io"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/raft"

	cc "github.com/vontikov/go-concurrent"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
)

// ErrUnknownCommand is the error returned by fsm.Apply if the command is
// unknown.
var ErrUnknownCommand = errors.New("unknown command")

type streamMap cc.Map
type mutexMap cc.Map

// FSM implements raft.FSM
type FSM struct {
	sync.Mutex
	logger  logging.Logger
	ctx     context.Context
	streams streamMap
	qs      queueMapPtr
	ds      dictMapPtr
	ms      mutexMap
	mo      sync.Once
}

type fsmCommand func(*FSM, *pb.ClusterCommand) interface{}

// NewFSM creates a new raft.FSM instance
func NewFSM(ctx context.Context) *FSM {
	return &FSM{
		logger:  logging.NewLogger("fsm"),
		ctx:     ctx,
		qs:      newQueueMap(),
		ds:      newDictMap(),
		ms:      cc.NewSynchronizedMap(0),
		streams: cc.NewSynchronizedMap(0),
	}
}

var fsmCommands [pb.ClusterCommand_MAX_INDEX]fsmCommand

func init() {
	fsmCommands = [...]fsmCommand{
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
}

// Apply log is invoked once a log entry is committed.
func (f *FSM) Apply(l *raft.Log) interface{} {
	var c pb.ClusterCommand
	if err := proto.Unmarshal(l.Data, &c); err != nil {
		return err
	}
	f.Lock()
	defer f.Unlock()
	return fsmCommands[c.Command](f, &c)
}

// Execute executes the command c.
func (f *FSM) Execute(c *pb.ClusterCommand) interface{} {
	f.Lock()
	defer f.Unlock()
	return fsmCommands[c.Command](f, c)
}

// Snapshot is used to support log compaction.
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return f, nil
}

// Persist implements raft.SnapshotSink.Persist.
func (f *FSM) Persist(sink raft.SnapshotSink) error {
	f.Lock()
	defer f.Unlock()

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
		sink.Cancel()
	}
	return err
}

// Restore is used to restore an FSM from a snapshot.
func (f *FSM) Restore(rc io.ReadCloser) error {
	f.Lock()
	defer f.Unlock()
	return gob.NewDecoder(rc).Decode(f)
}

// Release implements raft.SnapshotSink.Release.
func (f *FSM) Release() {}

// MarshalBinary implements encoding.BinaryMarshaler.MarshalBinary
// (https://golang.org/pkg/encoding/#BinaryMarshaler).
func (f *FSM) MarshalBinary() ([]byte, error) {
	f.logger.Debug("start marshalling")

	qb, err := f.qs.MarshalBinary()
	if err != nil {
		return nil, err
	}
	qbSz := len(qb)
	f.logger.Debug("queues marshalled", "size", qbSz)

	db, err := f.ds.MarshalBinary()
	if err != nil {
		return nil, err
	}
	dbSz := len(db)
	f.logger.Debug("dictionaries marshalled", "size", qbSz)

	data := make([]byte, 4+qbSz+4+dbSz)

	idx := 0

	// queues size
	binary.LittleEndian.PutUint32(data[idx:], uint32(qbSz))
	idx += 4
	// queues
	idx += copy(data[idx:], qb)

	// dictionaries size
	binary.LittleEndian.PutUint32(data[idx:], uint32(dbSz))
	idx += 4

	// dictionaries
	copy(data[idx:], db)

	f.logger.Debug("marshalling complete", "size", len(data))
	return data, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.UnmarshalBinary
// (https://golang.org/pkg/encoding/#BinaryUnmarshaler).
func (f *FSM) UnmarshalBinary(data []byte) error {
	f.logger.Debug("start unmarshalling", "size", len(data))
	idx := 0

	// queues size
	sz := int(binary.LittleEndian.Uint32(data[idx:]))
	idx += 4
	// queues
	qs := newQueueMap()
	if err := qs.UnmarshalBinary(data[idx : idx+sz]); err != nil {
		return err
	}
	idx += sz
	f.logger.Debug("queues unmarshalled", "size", sz)

	// dictionaries size
	sz = int(binary.LittleEndian.Uint32(data[idx:]))
	idx += 4
	// dictionaries
	ds := newDictMap()
	if err := ds.UnmarshalBinary(data[idx : idx+sz]); err != nil {
		return err
	}
	f.logger.Debug("dictionaries unmarshalled", "size", sz)

	f.Lock()
	f.qs = qs
	f.ds = ds
	f.Unlock()

	f.logger.Debug("unmarshalling complete")
	return nil
}
