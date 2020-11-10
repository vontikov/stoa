package cluster

import (
	"testing"

	"context"
	"sync"

	"github.com/stretchr/testify/assert"
)

func TestFsmDictionary(t *testing.T) {
	const dictName = "test-dictionary"
	const max = 100

	f := NewFSM(context.Background())

	var wg sync.WaitGroup
	for i := 0; i < max; i++ {
		wg.Add(1)
		go func(n int) {
			d := f.dictionary(dictName)
			assert.NotNil(t, d)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
