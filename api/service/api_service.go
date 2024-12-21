package service

import (
	"context"
	"fmt"

	"github.com/skip-mev/platform-take-home/api/types"
	"github.com/skip-mev/platform-take-home/observability/logging"
	"github.com/skip-mev/platform-take-home/store"
	"go.uber.org/zap"
)

type TakeHomeService struct {
	store *store.DBStore
	types.UnimplementedTakeHomeServiceServer
}

var _ types.TakeHomeServiceServer = &TakeHomeService{}

func NewTakeHomeService(store *store.DBStore) *TakeHomeService {
	return &TakeHomeService{store: store}
}

func (s *TakeHomeService) GetItems(ctx context.Context, _ *types.EmptyRequest) (*types.GetItemsResponse, error) {
	items, err := s.store.GetItems(ctx)

	if err != nil {
		logging.FromContext(ctx).Error("failed to retrieve items", zap.Error(err))
		return &types.GetItemsResponse{Items: make([]*types.Item, 0)}, fmt.Errorf("failed to retrieve items")
	}

	apiItems := make([]*types.Item, 0, len(items))

	for _, item := range items {
		apiItems = append(apiItems, &types.Item{
			Id:          uint64(item.ID),
			Name:        item.Name,
			Description: item.Description,
		})
	}

	return &types.GetItemsResponse{Items: apiItems}, nil
}

func (s *TakeHomeService) GetItem(ctx context.Context, req *types.GetItemRequest) (*types.GetItemResponse, error) {
	item, err := s.store.GetItem(ctx, uint(req.Id))

	if err != nil {
		logging.FromContext(ctx).Error("failed to retrieve item", zap.Error(err))
		return &types.GetItemResponse{}, fmt.Errorf("failed to retrieve item")
	}

	return &types.GetItemResponse{
		Item: &types.Item{
			Id:          uint64(item.ID),
			Name:        item.Name,
			Description: item.Description,
		},
	}, nil
}

func (s *TakeHomeService) CreateItem(ctx context.Context, req *types.CreateItemRequest) (*types.CreateItemResponse, error) {
	item, err := s.store.CreateItem(ctx, req.Item.Name, req.Item.Description)

	if err != nil {
		logging.FromContext(ctx).Error("failed to create item", zap.Error(err))
		return &types.CreateItemResponse{}, fmt.Errorf("failed to create item")
	}

	return &types.CreateItemResponse{ItemId: uint64(item)}, nil
}
