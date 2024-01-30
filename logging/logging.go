package logging

import (
	"context"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func DefaultLogger(options ...zap.Option) (*zap.Logger, error) {
	if os.Getenv("DEV_LOGGING") == "true" {
		return zap.NewDevelopment(options...)
	}

	return zap.NewProduction(options...)
}

func WithDefaultLogger(ctx context.Context, options ...zap.Option) (context.Context, error) {
	logger, err := DefaultLogger(options...)
	if err != nil {
		return ctx, err
	}
	return WithLogger(ctx, logger), nil
}

// prevent collisions with other packages
type key int

var loggerKey key

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}

	return logger
}

func UnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx := WithLogger(ctx, logger)
		return handler(newCtx, req)
	}
}
