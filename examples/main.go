package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vontikov/stoa/internal/logging"
	stoa "github.com/vontikov/stoa/pkg/client"
)

const dictName = "dict"

// Version contains the app version.
var Version string = "N/A"

var bootstrap = flag.String("bootstrap", "", "Comma separated peer list")
var logLevel = flag.String("log-level", "info", "Log level: trace|debug|info|warn|error|none")

func main() {
	flag.Parse()
	if *bootstrap == "" {
		panic("peer list must not be empty")
	}

	logging.SetLevel(*logLevel)
	logger := logging.NewLogger("stoa-example")

	logger.Info("starting", "version", Version)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		logger.Info("done")
	}()

	client, err := stoa.New(
		stoa.WithContext(ctx),
		stoa.WithPeers(*bootstrap),
	)
	panicOnError(err)
	logger.Info("client is ready")

	logger.Info("started")

	dict := client.Dictionary(dictName)
	logger.Info("dictionary created", "name", dictName)

	i := 0
	for {
		select {
		case sig := <-signals:
			logger.Debug("received signal", "type", sig)
			return
		default:
			k := fmt.Sprintf("k-%d", i)
			i++
			r, err := dict.PutIfAbsent(ctx, []byte(k), []byte("value"))
			panicOnError(err)
			logger.Info("dictionary putIfAbsent", "result", r)
			sz, err := dict.Size(ctx)
			panicOnError(err)
			logger.Info("dictionary size", "value", sz)
			time.Sleep(1 * time.Second)
		}
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
