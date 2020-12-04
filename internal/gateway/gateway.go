package gateway

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
)

// Gateway processes external requests.
type Gateway = errgroup.Group

// New creates new Gateway instance
func New(ctx context.Context, ip string, grpcPort int, httpPort int, cluster cluster.Cluster) (*Gateway, error) {
	logger := logging.NewLogger("gateway")

	grpcAddr := fmt.Sprintf("%s:%d", ip, grpcPort)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStoaServer(grpcServer, newServer(cluster))

	healthcheck := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthcheck)

	httpAddr := fmt.Sprintf("%s:%d", ip, httpPort)
	mux := runtime.NewServeMux()
	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	err = pb.RegisterStoaHandlerFromEndpoint(context.Background(), mux, grpcAddr, opts)
	if err != nil {
		return nil, err
	}

	// Prometheus metrics
	mux.Handle(
		"GET",
		runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"metrics"}, "")),
		func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			promhttp.Handler().ServeHTTP(w, r)
		})

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		logger.Info("serving gRPC", "address", grpcAddr)
		defer logger.Info("gRPC stopped")
		return grpcServer.Serve(lis)
	})
	g.Go(func() error {
		logger.Info("serving HTTP", "address", httpAddr)
		defer logger.Info("http stopped")
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-ctx.Done()
		httpServer.Close()
		grpcServer.GracefulStop()
		return nil
	})

	return g, nil
}
