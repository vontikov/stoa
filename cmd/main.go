package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/gateway"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/internal/metric"
	"github.com/vontikov/stoa/internal/util"
)

var (
	// App is the app name.
	App string = "stoa"
	// Version is the app version.
	Version string = "N/A"
)

var (
	grpcPort  = flag.Int("grpc-port", 3500, "gRPC port")
	httpPort  = flag.Int("http-port", 3501, "HTTP port")
	ip        = flag.String("ip", "0.0.0.0", "IP")
	logLevel  = flag.String("log-level", "info", "Log level: trace|debug|info|warn|error|none")
	bootstrap = flag.String("bootstrap", "", "Raft bootstrap")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	logging.SetLevel(*logLevel)
	logger := logging.NewLogger(App)
	logOptions(logger)
	logger.Info("starting", "version", Version)

	if *ip == "" {
		ifaces, err := util.GetInterfaces()
		panicOnError(err)
		if len(ifaces) == 0 {
			panic("networking interfaces not found")
		}
		*ip = ifaces[0].String()
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	var err error
	peers := *bootstrap
	if peers == "" {
		peers, err = os.Hostname()
		panicOnError(err)
	}
	cluster, err := cluster.New(cluster.WithPeers(peers))
	panicOnError(err)

	metric.Init(App, Version, cluster)

	ctx, cancel := context.WithCancel(context.Background())
	gateway, err := gateway.New(ctx, *ip, *grpcPort, *httpPort, cluster)
	panicOnError(err)

	metric.Info.Set(1.0)
	logger.Info("started")

	sig := <-signals
	logger.Debug("received signal", "type", sig)
	cancel()
	gateway.Wait()
	cluster.Shutdown()
	logger.Info("done")
}

func logOptions(l logging.Logger) {
	l.Debug("option", "gRPC port", *grpcPort)
	l.Debug("option", "HTTP port", *httpPort)
	l.Debug("option", "IP", *ip)
	l.Debug("option", "peers", *bootstrap)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
