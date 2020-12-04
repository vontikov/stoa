package gateway

import (
	"io/ioutil"
	"testing"

	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
)

const (
	testLogLevel = "debug"
)

func init() {
	logging.SetLevel(testLogLevel)
}

func RunTestCluster(ctx context.Context, t *testing.T) {
	const (
		addr1 = "127.0.0.1:3501"
		addr2 = "127.0.0.1:3502"
		addr3 = "127.0.0.1:3503"
	)

	peers := []struct {
		grpcPort int
		httpPort int
		peers    string
	}{

		{2501, 2601, strings.Join([]string{addr1, addr2, addr3}, cluster.PeerListSep)},
		{2502, 2602, strings.Join([]string{addr2}, cluster.PeerListSep)},
		{2503, 2603, strings.Join([]string{addr3}, cluster.PeerListSep)},
	}

	var wg sync.WaitGroup
	var m sync.Mutex
	var clusters []cluster.Cluster
	for _, p := range peers {
		wg.Add(1)
		go func(grpcPort, httpPort int, peers string) {
			cluster, err := cluster.New(cluster.WithPeers(peers))
			if err != nil {
				t.Logf("cluster error: %v", err)
				t.Fail()
			}
			m.Lock()
			clusters = append(clusters, cluster)
			m.Unlock()
			gateway, err := New(ctx, "localhost", grpcPort, httpPort, cluster)
			if err != nil {
				t.Logf("gateway error: %v", err)
				t.Fail()
			}
			wg.Done()

			<-ctx.Done()
			gateway.Wait()
			cluster.Shutdown()
		}(p.grpcPort, p.httpPort, p.peers)
	}
	wg.Wait()

	// make sure the leader is the same on all the peers
	timeout := time.After(10 * time.Second)
loop:
	for {
		select {
		case <-timeout:
			t.Logf("Election timeout")
			t.Fail()
		default:
			for _, c := range clusters {
				if c.LeaderAddress() != addr1 {
					continue loop
				}
			}
			break loop
		}
	}
}

func TestGateway(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	RunTestCluster(ctx, t)

	kv := pb.KeyValue{
		Name:  "test",
		Key:   []byte("foo"),
		Value: []byte("bar"),
	}

	const url = "http://127.0.0.1:%d/v1/dictionary/put"
	const contType = "application/json"

	var buf bytes.Buffer
	m := jsonpb.Marshaler{}

	// OK
	err := m.Marshal(&buf, &kv)
	assert.Nil(err)
	resp1, err := http.Post(fmt.Sprintf(url, 2601), contType, &buf)
	assert.Nil(err)
	defer resp1.Body.Close()
	assert.Equal(200, resp1.StatusCode)

	// valid request to a follower
	err = m.Marshal(&buf, &kv)
	assert.Nil(err)
	resp2, err := http.Post(fmt.Sprintf(url, 2602), contType, &buf)
	assert.Nil(err)
	defer resp2.Body.Close()
	assert.Equal(400, resp2.StatusCode)
	body, err := ioutil.ReadAll(resp2.Body)
	assert.Nil(err)
	assert.True(strings.Contains(string(body), raft.ErrNotLeader.Error()))

	// invalid request to a follower
	resp3, err := http.Post(fmt.Sprintf(url, 2603), contType, &buf)
	assert.Nil(err)
	defer resp3.Body.Close()
	assert.Equal(500, resp3.StatusCode)

	// invalid request to a follower
	resp1, err = http.Post(fmt.Sprintf(url, 2601), contType, &buf)
	assert.Nil(err)
	defer resp1.Body.Close()
	assert.Equal(500, resp1.StatusCode)
}
