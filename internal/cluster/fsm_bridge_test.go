package cluster

import (
	"testing"

	"context"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vontikov/stoa/pkg/pb"
)

func TestFSMPing(t *testing.T) {
	const (
		id1 = "id1"
		id2 = "id2"
	)

	te := time.Time{}
	tt := te
	timeNow = func() time.Time {
		tt = tt.Add(1 * time.Second)
		return tt
	}

	f := NewFSM(context.Background())

	mx := f.mutex("name")
	assert.Equal(t, te, mx.touched)

	r := mx.tryLock(id1)
	assert.True(t, r)
	assert.Equal(t, te.Add(1*time.Second), mx.touched)

	r = mx.tryLock(id1)
	assert.False(t, r)
	assert.Equal(t, te.Add(1*time.Second), mx.touched)

	f.ping(&pb.Ping{Id: id2})
	assert.Equal(t, te.Add(1*time.Second), mx.touched)

	f.ping(&pb.Ping{Id: id1})
	assert.Equal(t, te.Add(2*time.Second), mx.touched)
}
