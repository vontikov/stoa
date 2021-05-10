package main

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	stoa "github.com/vontikov/stoa/pkg/client"
)

const (
	bootstrap = "localhost:3001,localhost:3002,localhost:3003"
	mutexName = "mutex"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < 10; i++ {
		go f(ctx, strconv.Itoa(i))
	}

	<-signals
	cancel()
}

func f(ctx context.Context, id string) {
	client, err := stoa.New(ctx,
		stoa.WithID([]byte(id)),
		stoa.WithBootstrap(bootstrap),
	)
	if err != nil {
		log.Fatal("client error: ", err)
	}

	m := client.Mutex(mutexName)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			r, p, err := m.TryLock(ctx, []byte(id))
			if err != nil {
				log.Fatal("Lock: ", err)
			}
			if r {
				log.Printf("%s aquired the lock", id)

				time.Sleep(2000 * time.Millisecond)
				_, _, err = m.Unlock(ctx)
				if err != nil {
					log.Fatal("Unlock: ", err)
				}
				log.Printf("%s released the lock", id)
			} else {
				log.Printf("%s received payload from %s", id, string(p))
			}

			b, _ := rand.Int(rand.Reader, big.NewInt(1000))
			time.Sleep(time.Duration(b.Int64()) * time.Millisecond)
		}
	}
}
