package cluster

import (
	"context"
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

// entry is an data object stored in the FSM
type entry struct {
	Value     []byte // Value is the object's value
	TTLMillis int64  // TTLMillis is the object's TTL
}

// FSM implements raft.FSM
type FSM struct {
	logger  logging.Logger
	ctx     context.Context
	streams cc.Map
	qs      cc.Map
	ds      cc.Map
	ms      cc.Map
	mo      sync.Once
}

// NewFSM creates a new raft.FSM instance
func NewFSM(ctx context.Context) *FSM {
	return &FSM{
		logger:  logging.NewLogger("fsm"),
		ctx:     ctx,
		streams: cc.NewSynchronizedMap(0),
		qs:      cc.NewSynchronizedMap(0),
		ds:      cc.NewSynchronizedMap(0),
		ms:      cc.NewSynchronizedMap(0),
	}
}

// Apply log is invoked once a log entry is committed.
func (f *FSM) Apply(l *raft.Log) interface{} {
	var m pb.ClusterCommand
	if err := proto.Unmarshal(l.Data, &m); err != nil {
		return err
	}

	switch m.Command {
	case pb.ClusterCommand_QUEUE_SIZE:
		return f.queueSize(&m)
	case pb.ClusterCommand_QUEUE_CLEAR:
		return f.queueClear(&m)
	case pb.ClusterCommand_QUEUE_OFFER:
		return f.queueOffer(&m)
	case pb.ClusterCommand_QUEUE_POLL:
		return f.queuePoll(&m)
	case pb.ClusterCommand_QUEUE_PEEK:
		return f.queuePeek(&m)

	case pb.ClusterCommand_DICTIONARY_SIZE:
		return f.dictionarySize(&m)
	case pb.ClusterCommand_DICTIONARY_CLEAR:
		return f.dictionaryClear(&m)
	case pb.ClusterCommand_DICTIONARY_PUT:
		return f.dictionaryPut(&m)
	case pb.ClusterCommand_DICTIONARY_PUT_IF_ABSENT:
		return f.dictionaryPutIfAbsent(&m)
	case pb.ClusterCommand_DICTIONARY_GET:
		return f.dictionaryGet(&m)
	case pb.ClusterCommand_DICTIONARY_REMOVE:
		return f.dictionaryRemove(&m)
	case pb.ClusterCommand_DICTIONARY_SCAN:
		return f.dictionaryScan(&m)

	case pb.ClusterCommand_MUTEX_TRY_LOCK:
		return f.mutexTryLock(&m)
	case pb.ClusterCommand_MUTEX_UNLOCK:
		return f.mutexUnlock(&m)

	default:
		return ErrUnknownCommand
	}
}

// Snapshot is used to support log compaction.
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &snapshot{}, nil
}

// Restore is used to restore an FSM from a snapshot.
func (f *FSM) Restore(rc io.ReadCloser) error {
	panic("Restore()")
}
