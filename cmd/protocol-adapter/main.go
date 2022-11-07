package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"protocol-adapter/internal/config"
	"protocol-adapter/internal/grpcserver"
	"protocol-adapter/internal/httpserver"
	"protocol-adapter/internal/interceptors"
	"protocol-adapter/internal/logger"
)

func main() {
	// app settings
	appConfig := config.NewAppConf()
	// logger creation
	coreLog := logger.NewLogger(&appConfig.Meta, &appConfig.Logger)
	defer logger.Sync(coreLog)

	mainLog := coreLog.Named("main")
	defer logger.Sync(mainLog)

	logger.Info(context.TODO(), mainLog, "Application start")

	// metrics
	reg := prometheus.NewRegistry()
	grpcMetrics := grpc_prometheus.NewServerMetrics()
	reg.MustRegister(grpcMetrics)

	// Create HTTP server for prometheus.
	HTTPServerLog := coreLog.Named("HTTPServer")
	defer logger.Sync(HTTPServerLog)
	httpServer := httpserver.NewHTTPServer(reg, HTTPServerLog, appConfig.Metrics)

	// init httpclients
	// Client := smartclient.NewClient(
	// 	appConfig.Services.<Name>,
	// 	appConfig.Meta,
	// )

	// init services
	// Service := services.NewService(Client)

	// Prepare GRPC Server
	GRPCServerLog := coreLog.Named("GRPCServer")
	defer logger.Sync(GRPCServerLog)

	unaryInterceptor := []grpc.UnaryServerInterceptor{
		grpcMetrics.UnaryServerInterceptor(),
		interceptors.ZapLogUnaryServerInterceptor(GRPCServerLog),
		grpc_recovery.UnaryServerInterceptor(
			grpc_recovery.WithRecoveryHandlerContext(interceptors.PanicRecoverHandler),
		),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		grpcMetrics.StreamServerInterceptor(),
	}

	grpcServer := grpcserver.NewGRPCServer(GRPCServerLog, unaryInterceptor, streamInterceptors)
	// grpcServer.RegisterServices(grpcMetrics, Service)

	// RUN
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	// Cancel context by os.Signal
	go func() {
		defer wg.Done()
		sig := <-shutdownChannel
		logger.Error(
			context.TODO(), mainLog,
			"Received signal",
			zap.String("signal", sig.String()),
		)
		cancel()
	}()

	// Start http server for prometheus.
	wg.Add(1)
	go httpServer.Run(ctx, &wg)

	wg.Add(1)
	go grpcServer.Run(ctx, &wg, appConfig.GRPCServer.Port)

	wg.Wait()
	logger.Info(context.TODO(), mainLog, "Application has stopped")
}
