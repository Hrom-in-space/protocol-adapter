package grpcserver

import (
	"context"
	"fmt"
	"net"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"protocol-adapter/internal/logger"
)

type Servicer interface {
	ServiceDesc() *grpc.ServiceDesc
}

type GRPCServer struct {
	server *grpc.Server
	log    *zap.Logger
}

func NewGRPCServer(
	log *zap.Logger,
	ui []grpc.UnaryServerInterceptor,
	si []grpc.StreamServerInterceptor,
) *GRPCServer {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(ui...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(si...)),
	)

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	grpcServer := &GRPCServer{
		server: server,
		log:    log,
	}

	return grpcServer
}

func (s *GRPCServer) RegisterServices(grpcMetrics *grpc_prometheus.ServerMetrics, svcs ...Servicer) {
	for _, svc := range svcs {
		s.server.RegisterService(svc.ServiceDesc(), svc)
	}
	grpcMetrics.InitializeMetrics(s.server)
}

func (s *GRPCServer) Run(ctx context.Context, wg *sync.WaitGroup, port string) {
	defer wg.Done()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		logger.Fatal(context.TODO(), s.log, "Failed to listen", zap.Error(err))
	}

	// Start GRPC Server.
	go func() {
		logger.Info(context.TODO(), s.log, "GRPC server is starting", zap.Stringer("address", lis.Addr()))
		err = s.server.Serve(lis)
		logger.Info(context.TODO(), s.log, "GRPC server has stopped", zap.Error(err))
	}()

	s.gracefulStop(ctx)
}

func (s *GRPCServer) gracefulStop(ctx context.Context) {
	<-ctx.Done()
	s.server.GracefulStop()
	logger.Info(context.TODO(), s.log, "Graceful shutdown")
}
