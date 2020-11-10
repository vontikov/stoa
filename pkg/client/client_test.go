package client

import (
	"testing"

	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/discovery"
	"github.com/vontikov/stoa/internal/gateway"
	"github.com/vontikov/stoa/internal/logging"

	"github.com/vontikov/stoa/internal/util"
	"github.com/vontikov/stoa/pkg/pb"
)

const (
	testPeerID        = "id"
	testBindIP        = "127.0.0.1"
	testBindPort      = 3601
	testGrpcIP        = "127.0.0.1"
	testGrpcPort      = 3602
	testHTTPPort      = 3603
	testDiscoveryIP   = "224.0.0.1"
	testDiscoveryPort = 1234
	testLogLevel      = "debug"
)

func init() {
	logging.SetLevel(testLogLevel)
}

func runTestNode(t *testing.T) func() {
	ctx, cancel := context.WithCancel(context.Background())

	fsm := cluster.NewFSM(ctx)
	peer, err := cluster.NewPeer(testPeerID, testBindIP, testBindPort, fsm)
	assert.Nil(t, err)
	gw, err := gateway.New(testGrpcIP, testGrpcPort, testHTTPPort, peer, fsm)
	assert.Nil(t, err)

	ms := func() ([]byte, error) {
		return proto.Marshal(
			&pb.Discovery{
				Id:       testPeerID,
				GrpcIp:   testGrpcIP,
				GrpcPort: testGrpcPort,
				Leader:   peer.IsLeader(),
			})
	}
	sender, err := discovery.NewSender(testDiscoveryIP, testDiscoveryPort, ms)
	assert.Nil(t, err)

	go sender.Run(ctx)

	err = util.WaitLeaderStatus(peer, 5*time.Second)
	assert.Nil(t, err)

	return func() {
		gw.Shutdown()
		peer.Shutdown()
		cancel()
	}
}
