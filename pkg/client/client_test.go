package client

import (
	"sync"
	"testing"
	"time"

	"context"
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/internal/test"
)

func TestClientStart(t *testing.T) {
	const (
		clusterSize = 3
		basePort    = 2100
		dictName    = "dict"
	)

	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, bootstrap, err := test.StartCluster(ctx, basePort, clusterSize)
	assert.Nil(err)

	client, err := New(WithContext(ctx), WithPeers(bootstrap))
	assert.Nil(err)

	dict := client.Dictionary(dictName)
	sz, err := dict.Size(ctx)
	assert.Nil(err)
	assert.Equal(uint32(0), sz)
}

func TestClientWatch(t *testing.T) {
	const (
		clusterSize = 3
		basePort    = 2200
		queueName   = "test-queue"
		max         = 10000
		d           = max / 10
	)

	logging.SetLevel("debug")

	assert := assert.New(t)

	clusterCtx, clusterCancel := context.WithCancel(context.Background())
	peers, bootstrap, err := test.StartCluster(clusterCtx, basePort, clusterSize)
	assert.Nil(err)

	startProducer := make(chan StateChan)
	disrupt := make(chan bool)

	var wg sync.WaitGroup

	// disruptor
	disruptorCtx, disruptorCancel := context.WithCancel(context.Background())
	go func() {
		defer t.Log("disruptor complete")
		for {
			select {
			case <-disruptorCtx.Done():
			case <-disrupt:
				for _, p := range peers {
					if p.IsLeader() {
						t.Log("transfer leadership")
						p.LeadershipTransfer()
						break
					}
				}
			}
		}
	}()

	// consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer t.Log("consumer complete")

		clientCtx, clientCancel := context.WithCancel(context.Background())
		defer clientCancel()
		client, err := New(
			WithContext(clientCtx),
			WithPeers(bootstrap),
			WithLoggerName("stoa-consumer"),
		)
		assert.Nil(err)

		// in order to receive all the statuses, start producer when the client
		// is ready
		startProducer <- client.State()

		q := client.Queue(queueName)
		ch := q.Watch()

		timeout := time.After(60 * time.Second)
		n := 0
	loop:
		for {
			select {
			case <-timeout:
				t.Error("timeout")
				return
			default:
				m := <-ch
				n++
				if n%d == 0 {
					t.Logf("-> %d: %v", n, m)
				}
				if n == max {
					break loop
				}
			}
		}

		sz, err := q.Size(context.Background())
		assert.Nil(err)
		assert.Equal(uint32(max), sz)
	}()

	// producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer t.Log("producer complete")

		// make sure the receiver is ready
		state := <-startProducer

		clientCtx, clientCancel := context.WithCancel(context.Background())
		defer clientCancel()

		client, err := New(
			WithContext(clientCtx),
			WithPeers(bootstrap),
			WithLoggerName("stoa-producer"),
		)
		assert.Nil(err)

		q := client.Queue(queueName)

		for n := 0; n < max; n++ {
			msg := []byte(fmt.Sprintf("msg #%d", n))
			err := q.Offer(context.Background(), msg)
			assert.Nil(err)

			if n%d == 0 {
				t.Logf("%d ->", n)
				disrupt <- true

				// make sure the consumer is connected
				timeout := time.After(5 * time.Second)
			loop:
				for {
					select {
					case <-timeout:
						t.Error("timeout")
						return
					case s := <-state:
						if s == Connected {
							break loop
						}
					}
				}
			}
		}

		sz, err := q.Size(context.Background())
		assert.Nil(err)
		assert.Equal(uint32(max), sz)
	}()

	wg.Wait()
	disruptorCancel()
	clusterCancel()
}
