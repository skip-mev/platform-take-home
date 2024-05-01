package main

import (
	"context"
	"fmt"
	"net"
	"os"

	apiserver "github.com/skip-mev/platform-take-home/api/server"
	"github.com/skip-mev/platform-take-home/logging"
	"github.com/skip-mev/platform-take-home/types"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startGRPCServer(ctx context.Context, host string, port int) error {
	logger := logging.FromContext(ctx) // Use the logger extracted from context
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		logger.Error("Vault address is not set")
		return fmt.Errorf("vault address not set")
	}

	loggingInterceptor := logging.UnaryServerInterceptor(logger)
	server := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))

	// Create the API server instance with the vault address
	apiServer := apiserver.NewDefaultAPIServer(vaultAddr)
	types.RegisterAPIServer(server, apiServer)

	reflection.Register(server)

	go func() {
		<-ctx.Done()
		logger.Info("[grpc server] terminating...")
		server.GracefulStop()
	}()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		logger.Error("[grpc server] error creating listener", zap.Error(err))
		return fmt.Errorf("[grpc server] error creating listener: %v", err)
	}

	logger.Info("[grpc server] listening", zap.String("addr", fmt.Sprintf("http://%s", listener.Addr())))
	if err := server.Serve(listener); err != nil {
		logger.Error("[grpc server] error serving", zap.Error(err))
		return fmt.Errorf("[grpc server] error serving: %v", err)
	}

	logger.Info("[grpc server] terminated")
	return nil
}
