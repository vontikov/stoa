package gateway

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/logging"
	"github.com/vontikov/stoa/pkg/pb"
)

type Gateway struct {
	logger     logging.Logger
	wg         sync.WaitGroup
	lis        net.Listener
	restServer *http.Server
	grpcServer *grpc.Server
}

// New creates new Gateway instance
func New(ip string, grpcPort int, httpPort int, cluster cluster.Cluster) (*Gateway, error) {
	grpcAddr := fmt.Sprintf("%s:%d", ip, grpcPort)
	httpAddr := fmt.Sprintf("%s:%d", ip, httpPort)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStoaServer(grpcServer, newServer(cluster))

	healthcheck := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthcheck)

	mux := runtime.NewServeMux()
	restServer := &http.Server{
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

	g := Gateway{
		logger:     logging.NewLogger("gateway"),
		lis:        lis,
		restServer: restServer,
		grpcServer: grpcServer,
	}

	g.start()

	return &g, nil
}

func (g *Gateway) start() {
	g.wg.Add(1)
	go func() {
		g.grpcServer.Serve(g.lis)
		g.wg.Done()
	}()

	g.wg.Add(1)
	go func() {
		g.restServer.ListenAndServe()
		g.wg.Done()
	}()
}

func (g *Gateway) Shutdown() {
	g.restServer.Close()
	g.grpcServer.GracefulStop()
	g.wg.Wait()
	g.logger.Info("shutdown")
}
