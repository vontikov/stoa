package cluster

import (
	"testing"

	"bytes"
	"context"
	"encoding/gob"

	"github.com/stretchr/testify/assert"
)

func TestMarshaller(t *testing.T) {
	assert := assert.New(t)

	src := NewFSM(context.Background())
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(src)
	assert.Nil(err)

	dest := NewFSM(context.Background())
	err = gob.NewDecoder(&b).Decode(dest)
	assert.Nil(err)

	assert.Equal(src, dest)
}
