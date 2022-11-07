package interceptors

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"protocol-adapter/internal/logger"
	"protocol-adapter/internal/utils"
)

func ZapLogUnaryServerInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requestStartTime := time.Now()
		newCtx := logger.PutRequestStartTimeInCtx(ctx, requestStartTime)

		fields := []zap.Field{zap.String("handler-name", info.FullMethod)}

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			requestIDHeaders := md.Get("x-request-id")
			if len(requestIDHeaders) > 0 {
				fields = append(fields, zap.String("request-id", requestIDHeaders[0]))
			}
		}

		requestLogger := log.With(fields...).Named("rpc")
		defer logger.Sync(requestLogger)

		newCtx = ctxzap.ToContext(newCtx, requestLogger)

		return handler(newCtx, req)
	}
}

func PanicRecoverHandler(ctx context.Context, panic interface{}) (err error) {
	log := ctxzap.Extract(ctx)
	// на 5 уровней ниже по стеку расположена ошибка рантайм
	// если конечно функция вызывается через пакет grpc_recovery
	pc, filename, line, _ := runtime.Caller(5)
	filename = utils.PackagePath(filename)
	funcName := utils.RemoveBeforeFirst(runtime.FuncForPC(pc).Name(), ".")

	if err, ok := panic.(error); ok {
		logger.Info(
			ctx, log, "Runtime error",
			zap.Error(err),
			zap.String("src", fmt.Sprintf("%s:%v:%s", filename, line, funcName)),
		)
	} else {
		logger.Info(
			ctx, log, "Unknown runtime error",
			zap.String("panic", fmt.Sprintf("%T:%v", panic, panic)),
			zap.String("src", fmt.Sprintf("%s:%v:%s", filename, line, funcName)),
		)
	}

	return status.Error(codes.Internal, "Internal error")
}
