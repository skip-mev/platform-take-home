package main

import (
	"context"
	"github.com/skip-mev/platform-take-home/api/server"
	"github.com/skip-mev/platform-take-home/observability/logging"
	"github.com/skip-mev/platform-take-home/observability/metrics"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ctx, err := logging.WithDefaultLogger(ctx)

	if err != nil {
		panic(err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		grpcServer := server.NewServer()
		grpcServer.Start(ctx, "0.0.0.0:9008")
		return nil
	})

	eg.Go(func() error {
		if err := server.StartGRPCGateway(ctx, "0.0.0.0", 8080); err != nil {
			return err
		}

		return nil
	})

	eg.Go(func() error {
		if err := metrics.ServeMetrics(ctx, "0.0.0.0", 8081); err != nil {
			return err
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		logging.FromContext(ctx).Fatal("error during startup", zap.Error(err))
	}
}
