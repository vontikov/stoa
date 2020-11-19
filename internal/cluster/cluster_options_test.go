package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithPeersOption(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		peers []*Peer
		err   error
	}{
		{
			"Empty peer list",
			"",
			nil,
			ErrEmptyPeerList,
		},
		{
			"Empty peer list (trim)",
			"   ",
			nil,
			ErrEmptyPeerList,
		},
		{
			"Even peer list (1)",
			"id0:localhost:1234, id1:localhost:5678 ",
			nil,
			ErrEvenPeerList,
		},
		{
			"Even peer list (2)",
			"id0:localhost:1234, id1:localhost:5678,  id2:localhost:9101,  id3:localhost:9102",
			nil,
			ErrEvenPeerList,
		},
		{
			"Even peer list (3)",
			"id0:localhost:1234, id1:localhost:5678,  id2:localhost:9101,  id3:localhost:9103, ",
			nil,
			ErrEvenPeerList,
		},
		{
			"Incorrect peer parameters (1)",
			"id0:localhost:5678:xyz ",
			nil,
			ErrPeerParams,
		},
		{
			"Incorrect peer parameters (2)",
			"id0",
			nil,
			ErrPeerParams,
		},
		{
			"OK",
			"localhost:1234",
			[]*Peer{
				{"localhost:1234", "localhost", 1234},
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := newOptions(WithPeers(tt.input))
			assert.Equal(t, tt.err, err)
			if cfg != nil {
				assert.Equal(t, tt.peers, cfg.peers)
			}
		})
	}
}

func TestClusterOptions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		opts []Option
		err  error
	}{
		{
			"OK",
			[]Option{WithPeers("localhost:1234")},
			nil,
		},
		{
			"OK",
			nil,
			ErrClusterConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newOptions(tt.opts...)
			assert.Equal(t, tt.err, err)
		})
	}
}
