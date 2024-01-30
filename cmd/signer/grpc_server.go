package main

import (
	"context"
	"fmt"
	"net"

	apiserver "github.com/skip-mev/platform-take-home/api/server"
	"github.com/skip-mev/platform-take-home/logging"
	"github.com/skip-mev/platform-take-home/types"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startGRPCServer(ctx context.Context, host string, port int) error {
	loggingInterceptor := logging.UnaryServerInterceptor(logging.FromContext(ctx))

	server := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))

	types.RegisterAPIServer(server, apiserver.NewDefaultAPIServer(logging.FromContext(ctx)))

	reflection.Register(server)

	go func() {
		<-ctx.Done()
		logging.FromContext(ctx).Info("[grpc server] terminating...")
		server.GracefulStop()
	}()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return fmt.Errorf("[grpc server] error creating listener: %v", err)
	}

	logging.FromContext(ctx).Info("[grpc server] listening", zap.String("addr", fmt.Sprintf("http://%s", listener.Addr())))

	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("[grpc server] error serving: %v", err)
	}
	logging.FromContext(ctx).Info("[grpc server] terminated")

	return nil
}
