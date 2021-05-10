package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	stoa "github.com/vontikov/stoa/pkg/client"
)

const (
	bootstrap = "localhost:3001,localhost:3002,localhost:3003"
	queueName = "queue"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go offer(ctx)
	go poll(ctx, "1")
	go poll(ctx, "2")

	<-signals
	cancel()
}

func offer(ctx context.Context) {
	client, err := stoa.New(ctx, stoa.WithBootstrap(bootstrap))
	if err != nil {
		log.Fatal("client error: ", err)
	}

	q := client.Queue(queueName)

	t := time.NewTicker(1000 * time.Millisecond)
	defer t.Stop()

	n := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if err := q.Offer(ctx, []byte(strconv.Itoa(n))); err != nil {
				log.Fatal("offer: ", err)
			}
			log.Println("offer: ", n)
			n++
		}
	}
}

func poll(ctx context.Context, id string) {
	client, err := stoa.New(ctx,
		stoa.WithID([]byte(id)),
		stoa.WithBootstrap(bootstrap),
	)
	if err != nil {
		log.Fatal("client error: ", err)
	}

	q := client.Queue(queueName)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			v, err := q.Poll(ctx)
			if err != nil {
				log.Fatal("poll: ", err)
			}
			if v == nil {
				runtime.Gosched()
				break
			}

			log.Printf("client %s received: %s\n", id, string(v))
		}
	}
}
