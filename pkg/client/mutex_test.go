package client

import (
	"testing"

	"context"

	"github.com/stretchr/testify/assert"
	"github.com/vontikov/stoa/internal/test"
)

func TestMutex(t *testing.T) {
	const (
		clusterSize = 3
		basePort    = 2500
		muxName     = "mx"
	)

	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, peers, err := test.StartCluster(ctx, basePort, clusterSize)

	client0, err := New(ctx, WithBootstrap(peers))
	assert.Nil(err)

	mx0 := client0.Mutex(muxName)
	ok, _, err := mx0.TryLock(ctx, nil)
	assert.Nil(err)
	assert.True(ok)

	ok, _, err = mx0.TryLock(ctx, nil)
	assert.Nil(err)
	assert.False(ok)

	client1, err := New(ctx, WithBootstrap(peers))
	assert.Nil(err)

	mx1 := client1.Mutex(muxName)
	ok, _, err = mx1.TryLock(ctx, nil)
	assert.Nil(err)
	assert.False(ok)

	ok, _, err = mx0.Unlock(ctx)
	assert.Nil(err)
	assert.True(ok)
	ok, _, err = mx0.Unlock(ctx)
	assert.Nil(err)
	assert.False(ok)

	ok, _, err = mx1.Unlock(ctx)
	assert.Nil(err)
	assert.False(ok)
}
