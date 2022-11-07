package interceptors_test

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"protocol-adapter/internal/interceptors"
)

func TestPanicRecoverHandlerTable(t *testing.T) {
	type Panicer func()
	type Catcher func(context.Context, *observer.ObservedLogs)

	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	ctx := ctxzap.ToContext(context.Background(), observedLogger)

	tests := []struct {
		name      string
		panicFunc Panicer
		catchFunc Catcher
	}{
		{
			name: "Runtime error",
			catchFunc: func(ctx context.Context, obsLogs *observer.ObservedLogs) {
				p := recover()
				_ = interceptors.PanicRecoverHandler(ctx, p)

				require.Equal(t, 1, obsLogs.Len())
				allLogs := obsLogs.TakeAll()
				assert.Equal(t, "Runtime error", allLogs[0].Context[1].String)
			},
			panicFunc: func() {
				a := []int{}
				_ = a[1]
			},
		},
		{
			name: "Unknown runtime error",
			catchFunc: func(ctx context.Context, obsLogs *observer.ObservedLogs) {
				p := recover()
				_ = interceptors.PanicRecoverHandler(ctx, p)

				require.Equal(t, 1, obsLogs.Len())
				allLogs := obsLogs.TakeAll()
				assert.Equal(t, "Unknown runtime error", allLogs[0].Context[1].String)
				assert.Equal(t, "string:Anything", allLogs[0].Context[2].String)
			},
			panicFunc: func() {
				panic("Anything")
			},
		},
		{
			name: "Return internal status error",
			catchFunc: func(ctx context.Context, obsLogs *observer.ObservedLogs) {
				p := recover()
				err := interceptors.PanicRecoverHandler(ctx, p)
				assert.Equal(t, status.Code(err), codes.Internal)
			},
			panicFunc: func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.catchFunc(ctx, observedLogs)
			tt.panicFunc()
			observedLogs.TakeAll()
		})
	}
}
