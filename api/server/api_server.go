package server

import (
	"context"
	"github.com/skip-mev/platform-take-home/logging"
	"github.com/skip-mev/platform-take-home/types"
	"go.uber.org/zap"
)

type APIServerImpl struct {
	types.UnimplementedAPIServer

	logger *zap.Logger
}

var _ types.APIServer = &APIServerImpl{}

func NewDefaultAPIServer(logger *zap.Logger) *APIServerImpl {
	return &APIServerImpl{
		logger: logger,
	}
}

func (s *APIServerImpl) CreateWallet(ctx context.Context, request *types.CreateWalletRequest) (*types.CreateWalletResponse, error) {
	// TODO: implement this
	logging.FromContext(ctx).Info("CreateWallet", zap.String("name", request.Name))

	return &types.CreateWalletResponse{
		Wallet: &types.Wallet{},
	}, nil
}

func (s *APIServerImpl) GetWallet(ctx context.Context, request *types.GetWalletRequest) (*types.GetWalletResponse, error) {
	// TODO: implement this
	logging.FromContext(ctx).Info("GetWallet", zap.String("name", request.Name))

	return &types.GetWalletResponse{
		Wallet: &types.Wallet{},
	}, nil
}

func (s *APIServerImpl) GetWallets(ctx context.Context, request *types.EmptyRequest) (*types.GetWalletsResponse, error) {
	// TODO: implement this
	return &types.GetWalletsResponse{
		Wallets: nil,
	}, nil
}

func (s *APIServerImpl) Sign(ctx context.Context, request *types.WalletSignatureRequest) (*types.WalletSignatureResponse, error) {
	// TODO: implement this
	return &types.WalletSignatureResponse{
		Signature: nil,
	}, nil
}
