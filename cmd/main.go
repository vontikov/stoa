package main

import (
	"flag"
	"fmt"
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
	grpcPort      = flag.Int("grpc-port", 3502, "gRPC port")
	httpPort      = flag.Int("http-port", 3501, "HTTP port")
	logLevel      = flag.String("log-level", "info", "Log level: trace|debug|info|warn|error|none")
	ip            = flag.String("ip", "", "IP")
	peers         = flag.String("peers", "", "Peer addresses")
	pluginAddress = flag.String("plugin", "", "Plugin address")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	logging.SetLevel(*logLevel)
	logger := logging.NewLogger(App)
	logOptions(logger)
	logger.Info("Starting", "version", Version)

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
	logger.Info("Started")

	sig := <-signals
	logger.Debug("Received signal", "type", sig)
	gateway.Shutdown()
	cluster.Shutdown()
	logger.Info("Done")
}

func panicOnError(err error) {
	if err != nil {
		fmt.Println("=================")
		fmt.Println(err)
		fmt.Println("=================")
		panic(err)
	}
}

func logOptions(l logging.Logger) {
	l.Debug("option", "gRPC port", *grpcPort)
	l.Debug("option", "HTTP port", *httpPort)
	l.Debug("option", "IP", *ip)
	l.Debug("option", "peers", *peers)
	l.Debug("option", "pluginAddress", *pluginAddress)
}
