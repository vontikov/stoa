package client

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"bytes"
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/stretchr/testify/assert"
	"github.com/vontikov/stoa/internal/gateway"
	"github.com/vontikov/stoa/internal/test"
)

func TestDictionary(t *testing.T) {
	const (
		clusterSize = 3
		basePort    = 2300
		dictName    = "dict"
		max         = 100
	)

	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, bootstrap, err := test.StartCluster(ctx, basePort, clusterSize)

	client, err := New(ctx, WithBootstrap(bootstrap))
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

func TestDictionaryRange(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	const (
		clusterSize = 3
		basePort    = 2400
		dictName    = "dict"
		max         = 10000
	)

	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	peers, bootstrap, err := test.StartCluster(ctx, basePort, clusterSize)
	assert.Nil(err)

	client, err := New(ctx, WithBootstrap(bootstrap))
	assert.Nil(err)

	dict := client.Dictionary(dictName)

	// populate
	for i := 0; i < max; i++ {
		k := strconv.Itoa(i)
		v := strconv.Itoa(i + max)
		ok, err := dict.PutIfAbsent(ctx, []byte(k), []byte(v))
		assert.Nil(err)
		assert.True(ok)
	}

	tctx, tcancel := context.WithTimeout(ctx, 30*time.Second)
	defer tcancel()

	// transfer leadership
	var trs int32
	go func() {
		t := time.Tick(5 * time.Second)
		for {
			select {
			case <-tctx.Done():
				return
			case <-t:
				for _, c := range peers {
					if c.IsLeader() {
						err := c.LeadershipTransfer()
						assert.Nil(err)
						atomic.AddInt32(&trs, 1)
					}
				}
			}
		}
	}()

	var pos, neg int32

	testFunc := func() {
		m := make(map[int]int)
		kvChan, errChan := dict.Range(ctx)

	loop:
		for {
			select {
			case err := <-errChan:
				// happens while the new leader is being elected
				assert.True(errors.Is(err, gateway.ErrNotLeader))
				atomic.AddInt32(&neg, 1)
				return
			case kv := <-kvChan:
				if kv == nil {
					break loop
				}
				k, err := strconv.Atoi(string(kv[0]))
				assert.Nil(err)
				v, err := strconv.Atoi(string(kv[1]))
				assert.Nil(err)
				m[k] = v
			}
		}

		assert.Equal(max, len(m))
		for i := 0; i < max; i++ {
			assert.Equal(i+max, m[i])
		}

		atomic.AddInt32(&pos, 1)
	}

loop:
	for {
		select {
		case <-tctx.Done():
			break loop
		default:
			var wg sync.WaitGroup
			for i := 0; i < 3; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					testFunc()
					time.Sleep(1 * time.Second)
				}()
			}
			wg.Wait()
		}
	}

	t.Logf("trs=%d", trs)
	assert.Less(int32(3), trs)

	t.Logf("pos=%d", pos)
	assert.Less(int32(60), pos)

	t.Logf("neg=%d", neg)
	assert.Greater(int32(20), neg)
}
