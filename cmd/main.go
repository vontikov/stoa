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
	bindAddr   = flag.String("bind", "0.0.0.0", "Raft bind address")
	bootstrap  = flag.String("bootstrap", "", "Raft bootstrap")
	grpcPort   = flag.Int("grpc-port", 3500, "gRPC port")
	httpPort   = flag.Int("http-port", 3501, "HTTP port")
	listenAddr = flag.String("listen", "0.0.0.0", "Listen address")
	logLevel   = flag.String("log-level", "info", "Log level: trace|debug|info|warn|error|none")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()
	logging.SetLevel(*logLevel)
	logger := logging.NewLogger(App)

	hostname, _ := os.Hostname()
	logger.Info("starting", "name", App, "version", Version, "hostname", hostname)
	flag.VisitAll(func(f *flag.Flag) { logger.Debug("option", "name", f.Name, "value", f.Value) })

	if *listenAddr == "" {
		ifaces, err := util.GetInterfaces()
		panicOnError(err)
		if len(ifaces) == 0 {
			panic("networking interfaces not found")
		}
		*listenAddr = ifaces[0].String()
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	var err error
	peers := *bootstrap
	if peers == "" {
		peers, err = os.Hostname()
		panicOnError(err)
	}
	cluster, err := cluster.New(
		cluster.WithBindAddress(*bindAddr),
		cluster.WithPeers(peers),
	)
	panicOnError(err)

	metric.Init(App, Version, cluster)

	ctx, cancel := context.WithCancel(context.Background())
	gateway, err := gateway.New(ctx, *listenAddr, *grpcPort, *httpPort, cluster)
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

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
