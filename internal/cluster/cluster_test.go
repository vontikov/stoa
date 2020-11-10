package cluster

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/vontikov/stoa/internal/mocks/github.com/hashicorp/raft"
)

func TestClusterApply(t *testing.T) {
	const ip = "127.0.0.1"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	fsm := mock_raft.NewMockFSM(mockCtrl)
	fsm.EXPECT().Apply(gomock.Any())
	peer, err := NewPeer("id0", ip, 3500, fsm)
	assert.Nil(t, err)
	defer peer.Shutdown()

	fsm1 := mock_raft.NewMockFSM(mockCtrl)
	fsm1.EXPECT().Apply(gomock.Any())
	peer1, err := NewPeer("id1", ip, 3501, fsm1)
	assert.Nil(t, err)
	defer peer1.Shutdown()

	fsm2 := mock_raft.NewMockFSM(mockCtrl)
	fsm2.EXPECT().Apply(gomock.Any())
	peer2, err := NewPeer("id2", ip, 3502, fsm2)
	assert.Nil(t, err)
	defer peer2.Shutdown()

	// wait until the peer obtain leader status
	timeout := time.After(5 * time.Second)
loop:
	for {
		select {
		case <-timeout:
			t.Fail()
			t.Logf("Could not achieve leader status")
		default:
			if peer.IsLeader() {
				break loop
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	// add voters
	err = peer.AddVoter("id1", ip, 3501)
	assert.Nil(t, err)
	err = peer.AddVoter("id2", ip, 3502)
	assert.Nil(t, err)

	// apply a message to the log
	err = peer.R.Apply([]byte("hello"), 0).Error()
	assert.Nil(t, err)

	// wait until the message is propagated
	err = peer.R.Barrier(5 * time.Second).Error()
	assert.Nil(t, err)
}
