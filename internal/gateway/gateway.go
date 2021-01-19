package gateway

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
)

type options struct {
	grpcPort        int
	httpPort        int
	ip              string
	loggerName      string
	metricsEnabled  bool
	profilerEnabled bool
}

// Gateway processes external requests.
type Gateway = errgroup.Group

// Option defines a Gateway configuration option.
type Option func(*options)

func WithListenAddress(v string) Option { return func(o *options) { o.ip = v } }
func WithGRPCPort(v int) Option         { return func(o *options) { o.grpcPort = v } }
func WithHTTPPort(v int) Option         { return func(o *options) { o.httpPort = v } }
func WithMetricsEnabled(v bool) Option  { return func(o *options) { o.metricsEnabled = v } }
func WithPprofEnabled(v bool) Option    { return func(o *options) { o.profilerEnabled = v } }
func WithLoggerName(v string) Option    { return func(o *options) { o.loggerName = v } }

// New creates new Gateway instance
func New(ctx context.Context, cluster cluster.Cluster, opts ...Option) (*Gateway, error) {
	cfg := &options{}
	for _, o := range opts {
		o(cfg)
	}

	logger := logging.NewLogger(cfg.loggerName)

	grpcAddr := fmt.Sprintf("%s:%d", cfg.ip, cfg.grpcPort)
	httpAddr := fmt.Sprintf("%s:%d", cfg.ip, cfg.httpPort)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return nil, err
	}

	stoaServer := newServer(ctx, cluster, logger)
	grpcServer := grpc.NewServer()
	pb.RegisterStoaServer(grpcServer, stoaServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	gwMux := runtime.NewServeMux()
	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	err = pb.RegisterStoaHandlerFromEndpoint(context.Background(), gwMux, grpcAddr, dialOpts)
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()

	router.PathPrefix("/v1").Handler(gwMux)
	registerServiceHandlers(router, cfg)

	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: router,
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		logger.Info("serving grpc", "address", grpcAddr)
		defer logger.Info("grpc stopped")
		return grpcServer.Serve(lis)
	})
	g.Go(func() error {
		logger.Info("serving http", "address", httpAddr)
		defer logger.Info("http stopped")
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-ctx.Done()
		_ = httpServer.Close()
		grpcServer.GracefulStop()
		return nil
	})

	return g, nil
}

func registerServiceHandlers(router *mux.Router, cfg *options) {
	if cfg.metricsEnabled {
		router.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	}

	if cfg.profilerEnabled {
		pprofRouter := router.PathPrefix("/debug/pprof").Subrouter()
		pprofRouter.HandleFunc("", pprof.Index)
		pprofRouter.HandleFunc("/cmdline", pprof.Cmdline)
		pprofRouter.HandleFunc("/symbol", pprof.Symbol)
		pprofRouter.HandleFunc("/trace", pprof.Trace)

		profileRouter := pprofRouter.PathPrefix("/profile").Subrouter()
		profileRouter.HandleFunc("", pprof.Profile)
		profileRouter.Handle("/block", pprof.Handler("block"))
		profileRouter.Handle("/goroutine", pprof.Handler("goroutine"))
		profileRouter.Handle("/heap", pprof.Handler("heap"))
		profileRouter.Handle("/mutex", pprof.Handler("mutex"))
		profileRouter.Handle("/threadcreate", pprof.Handler("threadcreate"))
	}
}
