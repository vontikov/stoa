package client

import (
	"testing"

	"context"
	"fmt"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/test"
)

func TestClient(t *testing.T) {
	const (
		dictName = "dict"
		basePort = 2500
	)

	assert := assert.New(t)

	addr1 := fmt.Sprintf("127.0.0.1:%d", basePort+1)
	addr2 := fmt.Sprintf("127.0.0.1:%d", basePort+2)
	addr3 := fmt.Sprintf("127.0.0.1:%d", basePort+3)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	test.RunTestCluster(ctx, t, basePort)

	peers := strings.Join([]string{addr1, addr2, addr3}, cluster.PeerListSep)
	client, err := New(WithContext(ctx), WithPeers(peers))
	assert.Nil(err)

	dict := client.Dictionary(dictName)

	sz, err := dict.Size(ctx)
	assert.Nil(err)
	assert.Equal(uint32(0), sz)
}
