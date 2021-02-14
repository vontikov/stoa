package cluster

import (
	"testing"

	"context"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vontikov/stoa/pkg/pb"
)

func TestFSMPing(t *testing.T) {
	id1 := []byte("id1")
	id2 := []byte("id2")

	te := time.Time{}
	tt := te
	timeNow = func() time.Time {
		tt = tt.Add(1 * time.Second)
		return tt
	}

	f := NewFSM(context.Background())
	mx := mutex(f, "name")
	assert.Equal(t, te, mx.touched)

	r := mx.tryLock(id1)
	assert.True(t, r)
	assert.Equal(t, te.Add(1*time.Second), mx.touched)

	r = mx.tryLock(id1)
	assert.False(t, r)
	assert.Equal(t, te.Add(1*time.Second), mx.touched)

	processPing(f, &pb.ClusterCommand{
		Command: pb.ClusterCommand_SERVICE_PING,
		Payload: &pb.ClusterCommand_ClientId{ClientId: &pb.ClientId{Id: id2}},
	})
	assert.Equal(t, te.Add(1*time.Second), mx.touched)

	processPing(f, &pb.ClusterCommand{
		Command: pb.ClusterCommand_SERVICE_PING,
		Payload: &pb.ClusterCommand_ClientId{ClientId: &pb.ClientId{Id: id1}},
	})
	assert.Equal(t, te.Add(2*time.Second), mx.touched)
}
