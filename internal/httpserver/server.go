package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"protocol-adapter/internal/config"
	"protocol-adapter/internal/logger"
)

type httpServer struct {
	server *http.Server
	log    *zap.Logger
}

func NewHTTPServer(reg prometheus.Gatherer, log *zap.Logger, conf config.Metrics) *httpServer {
	server := &http.Server{
		Handler:      promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:         fmt.Sprintf("0.0.0.0:%s", conf.Port),
		WriteTimeout: conf.HTTPWriteTimeout,
		ReadTimeout:  conf.HTTPReadTimeout,
	}

	return &httpServer{
		server: server,
		log:    log,
	}
}

func (s *httpServer) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	go s.run()
	s.gracefulStop(ctx)
}

func (s *httpServer) run() {
	logger.Info(context.TODO(), s.log, "HTTP server is starting", zap.String("address", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil {
		logger.Error(context.TODO(), s.log, "HTTP server has stopped", zap.Error(err))
	}
}

// GraceStop - wait <-ctx.Done().
func (s *httpServer) gracefulStop(ctx context.Context) {
	<-ctx.Done()
	if err := s.server.Shutdown(ctx); err != nil {
		logger.Error(context.TODO(), s.log, "Graceful shutdown", zap.Error(err))
	}
}
