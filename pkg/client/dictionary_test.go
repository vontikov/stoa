package client

import (
	"testing"

	"bytes"
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/stretchr/testify/assert"
	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/test"
)

func TestDictionary(t *testing.T) {
	const (
		dictName = "dict"
		max      = 100
		basePort = 2600
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

	var wg sync.WaitGroup

	// PutIfAbsent, Put, Get
	for i := 0; i < max; i++ {
		wg.Add(1)

		go func(i int) {
			var k, v bytes.Buffer
			fmt.Fprintf(&k, "key-%d", i)
			fmt.Fprintf(&v, "value-%d", i)

			ok, err := dict.PutIfAbsent(ctx, k.Bytes(), v.Bytes())
			assert.Nil(err)
			assert.True(ok)

			o, err := dict.Put(ctx, k.Bytes(), v.Bytes())
			assert.Nil(err)
			assert.Equal(v.Bytes(), o)

			o, err = dict.Get(ctx, k.Bytes())
			assert.Nil(err)
			assert.Equal(v.Bytes(), o)

			wg.Done()
		}(i)
	}
	wg.Wait()
	sz, err = dict.Size(ctx)
	assert.Nil(err)
	assert.Equal(uint32(max), sz)

	// Scan
	var keys, vals []int
	rc, ec := dict.Range(ctx)
loop:
	for {
		select {
		case err := <-ec:
			t.Logf("%v", err)
			t.Fail()
		case kv := <-rc:
			if kv == nil {
				break loop
			}

			k, err := strconv.Atoi(strings.ReplaceAll(string(kv[0]), "key-", ""))
			assert.Nil(err)
			v, err := strconv.Atoi(strings.ReplaceAll(string(kv[1]), "value-", ""))
			assert.Nil(err)

			keys = append(keys, k)
			vals = append(vals, v)
		}
	}
	assert.Equal(max, len(keys))
	assert.Equal(max, len(vals))

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })

	for i := 0; i < max; i++ {
		assert.Equal(i, keys[i])
		assert.Equal(i, vals[i])
	}

	// Remove
	for i := 0; i < max; i++ {
		wg.Add(1)

		go func(i int) {
			var k bytes.Buffer
			fmt.Fprintf(&k, "key-%d", i)
			err := dict.Remove(ctx, k.Bytes())
			assert.Nil(err)
			wg.Done()
		}(i)
	}
	wg.Wait()
	sz, err = dict.Size(ctx)
	assert.Nil(err)
	assert.Equal(uint32(0), sz)
}
