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
	"github.com/vontikov/stoa/internal/metric"
	"github.com/vontikov/stoa/internal/util"
	"github.com/vontikov/stoa/pkg/logging"
)

var (
	// App is the app name.
	App string = "N/A"
	// Version is the app version.
	Version string = "N/A"
)

var (
	bindAddrFlag        = flag.String("bind", "0.0.0.0", "Raft bind address")
	bootstrapFlag       = flag.String("bootstrap", "", "Raft bootstrap")
	grpcPortFlag        = flag.Int("grpc-port", 3500, "gRPC port")
	httpPortFlag        = flag.Int("http-port", 3501, "HTTP port")
	listenAddrFlag      = flag.String("listen", "0.0.0.0", "Listen address")
	logLevelFlag        = flag.String("log-level", "info", "Log level: trace|debug|info|warn|error|none")
	metricsEnabledFlag  = flag.Bool("metrics", true, "Enable Prometheus metrics")
	profilerEnabledFlag = flag.Bool("profiler", false, "Enable profiler")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()
	logging.SetLevel(*logLevelFlag)
	logger := logging.NewLogger(App)

	ctx, cancel := context.WithCancel(context.Background())

	hostname, err := os.Hostname()
	util.PanicOnError(err)

	logger.Info("starting", "name", App, "version", Version, "hostname", hostname)
	flag.VisitAll(func(f *flag.Flag) { logger.Debug("option", "name", f.Name, "value", f.Value) })

	if *listenAddrFlag == "" {
		ifaces, err := util.GetInterfaces()
		util.PanicOnError(err)
		if len(ifaces) == 0 {
			panic("networking interfaces not found")
		}
		*listenAddrFlag = ifaces[0].String()
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	peers := *bootstrapFlag
	if peers == "" {
		peers = hostname
	}

	cluster, err := cluster.New(ctx,
		cluster.WithBindAddress(*bindAddrFlag),
		cluster.WithPeers(peers),
	)
	util.PanicOnError(err)

	metric.Init(App, Version, util.HostHostname(hostname), cluster)

	gateway, err := gateway.New(ctx, cluster,
		gateway.WithLoggerName(App),
		gateway.WithListenAddress(*listenAddrFlag),
		gateway.WithGRPCPort(*grpcPortFlag),
		gateway.WithHTTPPort(*httpPortFlag),
		gateway.WithMetricsEnabled(*metricsEnabledFlag),
		gateway.WithPprofEnabled(*profilerEnabledFlag),
	)
	util.PanicOnError(err)

	metric.Info.Set(1.0)
	logger.Info("started")

	sig := <-signals
	logger.Debug("received signal", "type", sig)
	cancel()
	_ = gateway.Wait()
	logger.Info("done")
}
