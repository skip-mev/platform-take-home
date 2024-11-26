package server

import (
	"context"
	"github.com/skip-mev/platform-take-home/api/service"
	"github.com/skip-mev/platform-take-home/observability/logging"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"os"

	"github.com/skip-mev/platform-take-home/api/types"
	"github.com/skip-mev/platform-take-home/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	grpcServer *grpc.Server
}

func NewServer() *Server {
	return &Server{
		grpcServer: grpc.NewServer(
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
		),
	}
}

func (s *Server) Start(ctx context.Context, address string) error {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		logging.FromContext(ctx).Fatal("error creating database connection", zap.Error(err))
		return err
	}

	var dbStore *store.DBStore

	if os.Getenv("POSTGRES_DSN") != "" {
		dbStore, err = store.NewPostgresBackedStore(os.Getenv("POSTGRES_DSN"))
		if err != nil {
			logging.FromContext(ctx).Fatal("error creating database connection", zap.Error(err))
			return err
		}
	} else {
		dbStore, err = store.NewSQLiteBackedStore()
		if err != nil {
			logging.FromContext(ctx).Fatal("error creating database connection", zap.Error(err))
			return err
		}
	}

	if err := dbStore.Migrate(); err != nil {
		logging.FromContext(ctx).Fatal("error migrating database", zap.Error(err))
		return err
	}

	takeHomeService := service.NewTakeHomeService(dbStore)

	types.RegisterTakeHomeServiceServer(s.grpcServer, takeHomeService)
	reflection.Register(s.grpcServer)

	go func() {
		<-ctx.Done()
		s.grpcServer.GracefulStop()
	}()

	if err := s.grpcServer.Serve(listener); err != nil {
		logging.FromContext(ctx).Fatal("error serving grpc", zap.Error(err))
		return err
	}

	return nil
}
