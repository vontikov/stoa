package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/vontikov/stoa/internal/logging"
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

	logging.SetLevel(*logLevel)
	logger := logging.NewLogger("stoa-example")

	logger.Info("starting", "version", Version)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	var client stoa.Client
	var err error

	if *bootstrap == "" {
		panic("bootstrap must not be empty")
	}
	client, err = stoa.New(
		stoa.WithContext(ctx),
		stoa.WithPeers(*bootstrap),
	)
	panicOnError(err)

	g, ctx := errgroup.WithContext(ctx)

	if *queueOutEnabled {
		g.Go(queueOut(ctx, client))
	}
	if *queueInEnabled {
		g.Go(queueIn(ctx, client))
	}
	logger.Info("started")

	sig := <-signals
	logger.Debug("received signal", "type", sig)
	cancel()
	g.Wait()
	logger.Info("done")
}

func queueOut(ctx context.Context, client stoa.Client) func() error {
	return func() error {
		logger := logging.NewLogger("queue-out")
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
				logger.Info("offer", "value", n)
				n++
			}
		}
	}
}

func queueIn(ctx context.Context, client stoa.Client) func() error {
	return func() error {
		logger := logging.NewLogger("queue-in")
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
				logger.Info("poll", "value", n)
			}
		}
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
