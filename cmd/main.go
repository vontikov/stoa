package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/gateway"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/internal/util"
)

var (
	// App contains the app name.
	App string = "stoa"
	// Version contains the app version.
	Version string = "N/A"
)

var (
	grpcPort = flag.Int("grpc-port", 3500, "gRPC port")
	httpPort = flag.Int("http-port", 3501, "HTTP port")
	ip       = flag.String("ip", "", "IP")
	logLevel = flag.String("log-level", "info", "Log level: trace|debug|info|warn|error|none")
	peers    = flag.String("peers", "", "Peer addresses")
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

	cluster, err := cluster.New(cluster.WithPeers(*peers))
	panicOnError(err)
	gateway, err := gateway.New(*ip, *grpcPort, *httpPort, cluster)
	panicOnError(err)
	logger.Info("started")

	sig := <-signals
	logger.Debug("received signal", "type", sig)
	gateway.Shutdown()
	cluster.Shutdown()
	logger.Info("done")
}

func logOptions(l logging.Logger) {
	l.Debug("option", "gRPC port", *grpcPort)
	l.Debug("option", "HTTP port", *httpPort)
	l.Debug("option", "IP", *ip)
	l.Debug("option", "peers", *peers)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
