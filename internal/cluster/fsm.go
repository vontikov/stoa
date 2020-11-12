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

type fsmCommand func(*FSM, *pb.ClusterCommand) interface{}

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
		pb.ClusterCommand_DICTIONARY_SCAN:          dictionaryScan,
		pb.ClusterCommand_MUTEX_TRY_LOCK:           mutexTryLock,
		pb.ClusterCommand_MUTEX_UNLOCK:             mutexUnlock,
	}
}

// Apply log is invoked once a log entry is committed.
func (f *FSM) Apply(l *raft.Log) interface{} {
	var c pb.ClusterCommand
	if err := proto.Unmarshal(l.Data, &c); err != nil {
		return err
	}
	return fsmCommands[c.Command](f, &c)
}

// Snapshot is used to support log compaction.
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &snapshot{}, nil
}

// Restore is used to restore an FSM from a snapshot.
func (f *FSM) Restore(rc io.ReadCloser) error {
	panic("Restore()")
}
