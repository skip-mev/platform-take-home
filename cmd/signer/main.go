package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/skip-mev/platform-take-home/logging"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	grpcHost        = "0.0.0.0"
	grpcPort        = 9000
	grpcGatewayHost = "0.0.0.0"
	grpcGatewayPort = 8080
)

func main() {
	loadEnv()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, err := logging.WithDefaultLogger(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating logger", err)
		os.Exit(1)
	}

	flag.Parse()

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := startGRPCServer(ctx, grpcHost, grpcPort); err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		if err := startGRPCGateway(ctx, grpcGatewayHost, grpcGatewayPort, fmt.Sprintf("%s:%d", grpcHost, grpcPort)); err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		logging.FromContext(ctx).Error("error running grpc service", zap.Error(err))
	}
}
