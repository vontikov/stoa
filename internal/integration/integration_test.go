package integration

/*
import (
	"testing"

	"context"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"

	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/discovery"
	"github.com/vontikov/stoa/internal/gateway"
	"github.com/vontikov/stoa/internal/util"
	"github.com/vontikov/stoa/pkg/client"
)

const (
	bindIP        = "127.0.0.1"
	grpcIP        = "127.0.0.1"
	discoveryIP   = "224.0.1.1"
	discoveryPort = 1334
)

func runTestNode(t *testing.T, id string, bindPort, grpcPort, httpPort int) (*cluster.Peer, func()) {
	assert := assert.New(t)
	g, ctx := errgroup.WithContext(context.Background())

	fsm := cluster.NewFSM(ctx)
	peer, err := cluster.NewPeer(id, bindIP, bindPort, fsm)
	assert.Nil(err)

	gw, err := gateway.New(grpcIP, grpcPort, httpPort, peer, fsm)
	assert.Nil(err)

	ms := discovery.DefaultMessageSupplier(peer, grpcIP, grpcPort)
	sender, err := discovery.NewSender(discoveryIP, discoveryPort, ms)
	assert.Nil(err)

	handler := discovery.DefaultHandler(peer, nil)
	receiver, err := discovery.NewReceiver(discoveryIP, discoveryPort, handler)
	assert.Nil(err)

	sender.Run(ctx, g)
	receiver.Run(ctx, g)

	return peer, func() {
		gw.Shutdown()
		peer.Shutdown()
	}
}

func TestIntegration(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	assert := assert.New(t)

	const (
		queueName = "queue"
		dictName  = "dict"
		muxName   = "mux"
		max       = 1000
	)

	defs := []struct {
		id       string
		bindPort int
		grpcPort int
		httpPort int
	}{
		{"id0", 3700, 3510, 3610},
		{"id1", 3701, 3511, 3611},
		{"id2", 3702, 3512, 3612},
		{"id3", 3703, 3513, 3613},
		{"id4", 3704, 3514, 3614},
	}

	var peers []*cluster.Peer
	var cancels []func()
	for _, d := range defs {
		p, c := runTestNode(t, d.id, d.bindPort, d.grpcPort, d.httpPort)
		peers = append(peers, p)
		cancels = append(cancels, c)
	}
	defer func() {
		for _, c := range cancels {
			c()
		}
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan interface{})

	// leadership transfer
	wg.Add(1)
	go func() {
		for {
			assert.Nil(util.WaitForSingleLeader(peers, 10*time.Second))
			idx := util.LeaderIndex(peers)
			assert.NotEqual(-1, idx)

			if v := <-sig; v == nil {
				break
			}

			err := peers[idx].LeadershipTransfer()
			assert.Nil(err)
		}
		wg.Done()
	}()

	// Queue.Poll()
	wg.Add(1)
	go func() {
		c, err := client.New(
			client.WithContext(ctx),
			client.WithDiscoveryIP(discoveryIP),
			client.WithDiscoveryPort(discoveryPort),
		)
		assert.Nil(err)

		q := c.Queue(queueName)
		n := 0
		for {
			qctx, qcancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer qcancel()
			b, err := q.Poll(qctx)
			assert.Nil(err)
			if b != nil {
				v, err := strconv.Atoi(string(b))
				assert.Nil(err)
				assert.Equal(n, v)
				n++
				if n == max {
					break
				}
				// transfer leadership
				if n%100 == 0 {
					sig <- n
				}
			}
		}
		close(sig)
		wg.Done()
	}()

	// Queue.Offer()
	wg.Add(1)
	go func() {
		c, err := client.New(
			client.WithContext(ctx),
			client.WithDiscoveryIP(discoveryIP),
			client.WithDiscoveryPort(discoveryPort),
		)
		assert.Nil(err)

		q := c.Queue(queueName)
		for i := 0; i < max; i++ {
			qctx, qcancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer qcancel()
			err := q.Offer(qctx, []byte(strconv.Itoa(i)))
			assert.Nil(err)
		}
		wg.Done()
	}()

	// Dictionary.Get()
	wg.Add(1)
	go func() {
		c, err := client.New(
			client.WithContext(ctx),
			client.WithDiscoveryIP(discoveryIP),
			client.WithDiscoveryPort(discoveryPort),
		)
		assert.Nil(err)

		d := c.Dictionary(dictName)
		for i := 0; i < max; i++ {
			tt := time.After(30 * time.Second)
			k := []byte(strconv.Itoa(i))
		loop:
			for {
				select {
				case <-tt:
					t.Logf("Timeout: %d", i)
					t.Fail()
				default:
					v, err := d.Get(ctx, k)
					assert.Nil(err)
					if v != nil {
						assert.Equal(k, v)
						break loop
					}
				}
			}
		}
		wg.Done()
	}()

	// Dictionary.Put()
	wg.Add(1)
	go func() {
		c, err := client.New(
			client.WithContext(ctx),
			client.WithDiscoveryIP(discoveryIP),
			client.WithDiscoveryPort(discoveryPort),
		)
		assert.Nil(err)

		d := c.Dictionary(dictName)
		for i := 0; i < max; i++ {
			v := []byte(strconv.Itoa(i))
			_, err := d.Put(ctx, v, v)
			assert.Nil(err)
		}
		wg.Done()
	}()

	// Mutex.Watch()
	wg.Add(1)
	go func() {
		c, err := client.New(
			client.WithContext(ctx),
			client.WithDiscoveryIP(discoveryIP),
			client.WithDiscoveryPort(discoveryPort),
		)
		assert.Nil(err)

		m := c.Mutex(muxName)
		w := struct{ client.MutexWatchProto }{}

		wc := make(chan interface{})
		var counter int32
		var value int32
		w.Callback = func(n string, v bool) {
			assert.Equal(muxName, n)

			var ov, nv int32
			if v {
				nv = 1
			} else {
				ov = 1
			}
			assert.True(atomic.CompareAndSwapInt32(&value, ov, nv))

			r := atomic.AddInt32(&counter, 1)
			if r == max*2 {
				wc <- true
			}
		}

		m.Watch(w)

		<-wc
		wg.Done()
	}()

	// Mutex.Lock()
	wg.Add(1)
	go func() {
		c, err := client.New(
			client.WithContext(ctx),
			client.WithDiscoveryIP(discoveryIP),
			client.WithDiscoveryPort(discoveryPort),
		)
		assert.Nil(err)

		m := c.Mutex(muxName)
		for i := 0; i < max; i++ {
			mctx, mcancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer mcancel()

			ok, err := m.TryLock(mctx)
			assert.Nil(err)
			assert.True(ok)

			time.Sleep(10 * time.Millisecond)

			ok, err = m.Unlock(mctx)
			assert.Nil(err)
			assert.True(ok)
		}
		wg.Done()
	}()

	wg.Wait()
	cancel()
}
*/
