package gateway

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/gogo/protobuf/jsonpb"
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

func startCluster(ctx context.Context, basePort int, peerNum int) (peers []cluster.Cluster, bootstrap string, err error) {
	if peerNum <= 0 {
		panic("peer number must be greater then 0")
	}
	if peerNum > 9 {
		panic("peer number must not be greater then 9")
	}
	if peerNum%2 == 0 {
		panic("peer number must not be even")
	}

	const (
		grpcPortShift = 1
		httpPortShift = 21
		bindPortShift = 11
	)

	type peer struct {
		grpcPort int
		httpPort int
		peers    string
	}

	var bindAddrs, bsAddrs []string
	var peerList []peer

	for i := 0; i < peerNum; i++ {
		bsAddrs = append(bsAddrs, fmt.Sprintf("127.0.0.1:%d", basePort+grpcPortShift+i))
		bindAddrs = append(bindAddrs, fmt.Sprintf("127.0.0.1:%d", basePort+bindPortShift+i))
	}

	for i := 0; i < peerNum; i++ {
		b := bindAddrs[i]
		if i == 0 {
			b = strings.Join(bindAddrs, cluster.PeerListSep)
		}
		peerList = append(peerList, peer{basePort + grpcPortShift + i, basePort + i + httpPortShift + i, b})
	}

	var m sync.Mutex
	var g errgroup.Group

	for i, p := range peerList {
		idx := i
		p := p
		g.Go(func() error {
			cluster, err := cluster.New(ctx,
				cluster.WithPeers(p.peers),
				cluster.WithLoggerName(fmt.Sprintf("raft-%d", idx)),
			)
			if err != nil {
				return err
			}
			m.Lock()
			peers = append(peers, cluster)
			m.Unlock()
			gateway, err := New(ctx, cluster,
				WithListenAddress("0.0.0.0"),
				WithGRPCPort(p.grpcPort),
				WithHTTPPort(p.httpPort),
				WithMetricsEnabled(true),
				WithPprofEnabled(true))
			if err != nil {
				return err
			}

			go func() {
				<-ctx.Done()
				_ = gateway.Wait()
			}()

			return nil
		})
	}

	err = g.Wait()
	if err != nil {
		return
	}

	// make sure the leader is the same on all the peers
	m.Lock()
	defer m.Unlock()
	leader, err := getLeader(peers, 10*time.Second)
	if err != nil {
		return
	}
	if leader != bindAddrs[0] {
		err = errors.New("unexpected leader address")
	}

	bootstrap = strings.Join(bsAddrs, cluster.PeerListSep)
	return
}

func getLeader(clusters []cluster.Cluster, timeout time.Duration) (leader string, err error) {
	t := time.After(timeout)
	for {
		select {
		case <-t:
			err = errors.New("election timeout")
			return
		default:
			r := true
			for _, c := range clusters {
				l := c.LeaderAddress()
				if l == "" {
					r = false
					break
				}
				if l != leader {
					leader = l
					r = false
					break
				}
			}
			if r {
				return
			}
		}
	}
}

func TestGatewayHTTP(t *testing.T) {
	const (
		clusterSize = 3
		basePort    = 3100
		dictName    = "dict"

		leaderPort   = basePort + 21
		followerPort = basePort + 23
	)

	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, _, err := startCluster(ctx, basePort, clusterSize)
	assert.Nil(err)

	t.Run("metric endpoint ", func(t *testing.T) {
		// leader
		func() {
			url := fmt.Sprintf("http://127.0.0.1:%d/metrics", leaderPort)
			resp, err := http.Get(url)
			assert.Nil(err)
			defer resp.Body.Close()
			assert.Equal(http.StatusOK, resp.StatusCode)
		}()

		// follower
		func() {
			url := fmt.Sprintf("http://127.0.0.1:%d/metrics", followerPort)
			resp, err := http.Get(url)
			assert.Nil(err)
			defer resp.Body.Close()
			assert.Equal(http.StatusOK, resp.StatusCode)
		}()
	})

	t.Run("leader dictionary", func(t *testing.T) {
		key := []byte("foo")
		value := []byte("bar")

		func() {
			kv := pb.KeyValue{
				Name:  dictName,
				Key:   key,
				Value: value,
			}

			url := fmt.Sprintf("http://127.0.0.1:%d/v1/dictionary/put", leaderPort)

			var buf bytes.Buffer
			m := jsonpb.Marshaler{}
			err := m.Marshal(&buf, &kv)
			assert.Nil(err)

			resp, err := http.Post(url, "application/json", &buf)
			assert.Nil(err)
			defer resp.Body.Close()
			assert.Equal(http.StatusOK, resp.StatusCode)
		}()

		func() {
			url := fmt.Sprintf("http://127.0.0.1:%d/v1/dictionary/get/%s/%s", leaderPort,
				dictName, url.QueryEscape(base64.StdEncoding.EncodeToString(key)))

			resp, err := http.Get(url)
			assert.Nil(err)
			defer resp.Body.Close()
			assert.Equal(http.StatusOK, resp.StatusCode)

			body, err := ioutil.ReadAll(resp.Body)
			assert.Nil(err)

			v := pb.Value{}
			m := jsonpb.Unmarshaler{}
			err = m.Unmarshal(bytes.NewBuffer(body), &v)
			assert.Nil(err)
			assert.Equal(value, v.Value)
		}()
	})

	t.Run("follower dictionary", func(t *testing.T) {
		key := []byte("foo")
		value := []byte("bar")

		func() {
			kv := pb.KeyValue{
				Name:  dictName,
				Key:   key,
				Value: value,
			}

			url := fmt.Sprintf("http://127.0.0.1:%d/v1/dictionary/put", followerPort)

			var buf bytes.Buffer
			m := jsonpb.Marshaler{}
			err := m.Marshal(&buf, &kv)
			assert.Nil(err)

			resp, err := http.Post(url, "application/json", &buf)
			assert.Nil(err)
			defer resp.Body.Close()
			assert.Equal(http.StatusBadRequest, resp.StatusCode)

			body, err := ioutil.ReadAll(resp.Body)
			assert.Nil(err)

			type rd struct {
				Code    int
				Message string
			}

			var actualBody rd
			err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&actualBody)
			assert.Nil(err)
			assert.Equal(rd{9, "node is not the leader"}, actualBody)
		}()
	})
}
