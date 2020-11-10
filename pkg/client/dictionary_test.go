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
)

func TestDictionary(t *testing.T) {
	const (
		dictName = "dict"
		max      = 100
	)

	defer runTestNode(t)()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := New(
		WithContext(ctx),
		WithDiscoveryIP(testDiscoveryIP),
		WithDiscoveryPort(testDiscoveryPort),
	)
	assert.Nil(t, err)
	assert.NotNil(t, client)

	dict := client.Dictionary(dictName)

	sz, err := dict.Size(ctx)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), sz)

	var wg sync.WaitGroup

	// PutIfAbsent, Put, Get
	for i := 0; i < max; i++ {
		wg.Add(1)

		go func(i int) {
			var k, v bytes.Buffer
			fmt.Fprintf(&k, "key-%d", i)
			fmt.Fprintf(&v, "value-%d", i)

			ok, err := dict.PutIfAbsent(ctx, k.Bytes(), v.Bytes())
			assert.Nil(t, err)
			assert.True(t, ok)

			o, err := dict.Put(ctx, k.Bytes(), v.Bytes())
			assert.Nil(t, err)
			assert.Equal(t, v.Bytes(), o)

			o, err = dict.Get(ctx, k.Bytes())
			assert.Nil(t, err)
			assert.Equal(t, v.Bytes(), o)

			wg.Done()
		}(i)
	}
	wg.Wait()
	sz, err = dict.Size(ctx)
	assert.Nil(t, err)
	assert.Equal(t, uint32(max), sz)

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
			assert.Nil(t, err)
			v, err := strconv.Atoi(strings.ReplaceAll(string(kv[1]), "value-", ""))
			assert.Nil(t, err)

			keys = append(keys, k)
			vals = append(vals, v)
		}
	}
	assert.Equal(t, max, len(keys))
	assert.Equal(t, max, len(vals))

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })

	for i := 0; i < max; i++ {
		assert.Equal(t, i, keys[i])
		assert.Equal(t, i, vals[i])
	}

	// Remove
	for i := 0; i < max; i++ {
		wg.Add(1)

		go func(i int) {
			var k bytes.Buffer
			fmt.Fprintf(&k, "key-%d", i)
			err := dict.Remove(ctx, k.Bytes())
			assert.Nil(t, err)
			wg.Done()
		}(i)
	}
	wg.Wait()
	sz, err = dict.Size(ctx)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), sz)
}
