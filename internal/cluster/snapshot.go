package cluster

import (
	"github.com/hashicorp/raft"
)

type snapshot struct {
}

func (*snapshot) Persist(sink raft.SnapshotSink) error {
	panic("Persist")
}

func (*snapshot) Release() {
}
