package main

import (
	"context"
	"fmt"
	"github.com/skip-mev/platform-take-home/logging"
	"github.com/skip-mev/platform-take-home/types"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startGRPCGateway(ctx context.Context, host string, port int, grpcServerEndpoint string) error {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	mux := runtime.NewServeMux()

	if err := types.RegisterAPIHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return fmt.Errorf("[grpc gateway] error creating listener: %v", err)
	}

	corsMiddleware := cors.New(cors.Options{})
	rootHandler := corsMiddleware.Handler(mux)

	server := &http.Server{Handler: rootHandler}

	// allows all origins by default
	go func() {
		<-ctx.Done()
		logging.FromContext(ctx).Info("[grpc gateway] terminating...")
		if err := server.Shutdown(context.Background()); err != nil {
			logging.FromContext(ctx).Error("[grpc gateway] failed to terminate", zap.Error(err))
		}
	}()

	logging.FromContext(ctx).Info("[grpc gateway] listening", zap.String("addr", fmt.Sprintf("http://%s", listener.Addr())))

	if err := server.Serve(listener); err != http.ErrServerClosed {
		return fmt.Errorf("[grpc gateway] error serving: %v", err)
	}
	logging.FromContext(ctx).Info("[grpc gateway] terminated")

	return nil
}
