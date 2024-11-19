package logging

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"

	"go.opentelemetry.io/otel/trace"
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

var (
	loggerKey       key = 0
	serviceLabelKey key = 1
)

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func WithServiceLabel(ctx context.Context, service string) context.Context {
	return context.WithValue(ctx, serviceLabelKey, service)
}

func TraceIDFromContext(ctx context.Context) (trace.TraceID, bool) {
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		return spanContext.TraceID(), true
	}
	return trace.TraceID{}, false
}

func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		var err error
		logger, err = DefaultLogger()
		if err != nil {
			return zap.NewNop()
		}
		logger.Error("missing logger on ctx", zap.Any("ctx", ctx))
	}

	traceID, ok := TraceIDFromContext(ctx)
	if ok {
		logger = logger.With(zap.String("traceId", traceID.String()))
	}

	service, ok := ctx.Value(serviceLabelKey).(string)
	if ok {
		logger = logger.With(zap.String("service", service))
	}

	return logger
}

func UnaryServerInterceptor(logger *zap.Logger, sampleRate float64) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		loggerForContext := logger
		traceID, ok := TraceIDFromContext(ctx)
		if ok {
			traceInt := binary.BigEndian.Uint16(traceID[:])
			if traceInt%100 > uint16(sampleRate*100) {
				loggerForContext = zap.NewNop()
			}
		}
		newCtx := WithLogger(ctx, loggerForContext)
		return handler(newCtx, req)
	}
}

func DefaultLoggingContext() context.Context {
	ctx, err := WithDefaultLogger(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating logger", err)
	}
	return ctx
}
