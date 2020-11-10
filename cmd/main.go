package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/cluster/plugin"
	"github.com/vontikov/stoa/internal/discovery"
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
	ip            = flag.String("ip", "", "IP")
	grpcPort      = flag.Int("grpc-port", 3502, "gRPC port")
	httpPort      = flag.Int("http-port", 3503, "HTTP port")
	discoveryIP   = flag.String("discovery-ip", "224.5.1.1", "Auto-discovery IP")
	discoveryPort = flag.Int("discovery-port", 3500, "Auto-discovery port")
	bindIP        = flag.String("bind-ip", "", "Bind IP")
	bindPort      = flag.Int("bind-port", 3501, "Bind port")
	pluginAddress = flag.String("plugin", "", "Plugin address")
	logLevel      = flag.String("log-level", "info", "Log level: trace|debug|info|warn|error|none")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	logging.SetLevel(*logLevel)
	logger := logging.NewLogger(App)

	logger.Debug("option", "IP", *ip)
	logger.Debug("option", "gRPC port", *grpcPort)
	logger.Debug("option", "HTTP port", *httpPort)
	logger.Debug("option", "Discovery IP", *discoveryIP)
	logger.Debug("option", "Discovery port", *discoveryPort)
	logger.Debug("option", "Bind IP", *bindIP)
	logger.Debug("option", "Bind port", *bindPort)

	logger.Info("Starting", "version", Version)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ifaces, err := util.GetInterfaces()
	panicOnError(err)
	if len(ifaces) == 0 {
		panic("networking interfaces not found")
	}
	addr := ifaces[0].String()

	if *bindIP == "" {
		bindIP = &addr
	}
	if *ip == "" {
		ip = &addr
	}

	ctx, cancel := context.WithCancel(context.Background())

	fsm := cluster.NewFSM(ctx)
	peerID := fmt.Sprintf("%s-%d", *bindIP, *bindPort)
	peer, err := cluster.NewPeer(peerID, *bindIP, *bindPort, fsm)
	panicOnError(err)

	handler := discovery.DefaultHandler(peer,
		func(args ...interface{}) {
			logger.Info("Discovered new peer", "address", args[0])
		})
	r, err := discovery.NewReceiver(*discoveryIP, *discoveryPort, handler)
	panicOnError(err)

	msgSuppl := discovery.DefaultMessageSupplier(peer, *ip, *grpcPort)
	s, err := discovery.NewSender(*discoveryIP, *discoveryPort, msgSuppl)
	panicOnError(err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("Receiver exited with error", "message", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("Sender exited with error", "message", err)
		}
	}()
	if *pluginAddress != "" {
		p := plugin.New(*pluginAddress, fsm)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := p.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
				logger.Error("Plugin exited with error", "message", err)
			}
		}()
	}

	gw, err := gateway.New(*ip, *grpcPort, *httpPort, peer, nil)
	panicOnError(err)
	logger.Info("Started")

	sig := <-signals
	logger.Debug("Received signal", "type", sig)
	cancel()
	gw.Shutdown()
	peer.Shutdown()
	logger.Info("Done")
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
