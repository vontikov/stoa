package cluster

import (
	"testing"

	"bytes"
	"encoding/gob"

	"github.com/stretchr/testify/assert"
)

func TestFSMEntryMarshalling(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		value []byte
		ttl   int64
	}{
		{nil, 0},
		{nil, 100},
		{[]byte{1, 2, 3, 4, 5}, 2000},
		{[]byte("Hello, World!"), 30000},
	}

	for _, tt := range tests {
		src := entry{
			Value:     tt.value,
			TTLMillis: tt.ttl,
		}

		var b bytes.Buffer
		err := gob.NewEncoder(&b).Encode(&src)
		assert.Nil(err)

		dst := entry{}
		err = gob.NewDecoder(&b).Decode(&dst)
		assert.Nil(err)
		assert.Equal(src, dst)
	}
}
