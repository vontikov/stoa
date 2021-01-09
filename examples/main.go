package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	stoa "github.com/vontikov/stoa/pkg/client"
)

const (
	dictName  = "dict"
	queueName = "queue"
)

// Version contains the app version.
var Version string = "N/A"

var (
	bootstrap       = flag.String("bootstrap", "", "Comma separated peer list")
	logLevel        = flag.String("log-level", "info", "Log level: trace|debug|info|warn|error|none")
	queueOutEnabled = flag.Bool("queue-out", false, "Enable output to the queue")
	queueInEnabled  = flag.Bool("queue-in", false, "Enable input from the queue")
)

func main() {
	flag.Parse()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	var client stoa.Client
	var err error

	if *bootstrap == "" {
		log.Fatal("bootstrap must not be empty")
	}
	client, err = stoa.New(
		stoa.WithContext(ctx),
		stoa.WithPeers(*bootstrap),
	)
	if err != nil {
		log.Fatal("client error: ", err)
	}

	g, ctx := errgroup.WithContext(ctx)

	if *queueOutEnabled {
		g.Go(queueOut(ctx, client))
	}
	if *queueInEnabled {
		g.Go(queueIn(ctx, client))
	}
	log.Println("started")

	<-signals
	cancel()
	g.Wait()
}

func queueOut(ctx context.Context, client stoa.Client) func() error {
	return func() error {
		q := client.Queue(queueName)

		t := time.Tick(1000 * time.Millisecond)
		n := 0
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-t:
				if err := q.Offer(ctx, []byte(strconv.Itoa(n))); err != nil {
					return err
				}
				log.Println("offer: ", n)
				n++
			}
		}
	}
}

func queueIn(ctx context.Context, client stoa.Client) func() error {
	return func() error {
		q := client.Queue(queueName)

		t := time.Tick(500 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-t:
				v, err := q.Poll(ctx)
				if err != nil {
					return err
				}
				if v == nil {
					runtime.Gosched()
					continue
				}

				n, err := strconv.Atoi(string(v))
				if err != nil {
					return err
				}
				log.Println("poll: ", n)
			}
		}
	}
}
