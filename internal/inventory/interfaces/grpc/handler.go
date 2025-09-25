package grpc_handler

import (
	"context"

	inventoryv1 "github.com/smallbiznis/go-genproto/smallbiznis/inventory"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryHandler struct {
	inventoryv1.UnimplementedInventoryServiceServer
}

type Params struct {
	fx.In
}

func NewInventoryHandler(p Params) *InventoryHandler {
	return &InventoryHandler{}
}

func (h *InventoryHandler) AdjustInventory(ctx context.Context, req *inventoryv1.AdjustInventoryRequest) (*inventoryv1.AdjustInventoryResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method AdjustInventory not implemented")
}

func (h *InventoryHandler) CreateVariantInventory(ctx context.Context, req *inventoryv1.InventoryItem) (*inventoryv1.InventoryItem, error) {
	return nil, status.Error(codes.Unimplemented, "method CreateVariantInventory not implemented")
}

func (h *InventoryHandler) GetVariantInventory(ctx context.Context, req *inventoryv1.GetVariantInventoryRequest) (*inventoryv1.GetVariantInventoryResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetVariantInventory not implemented")
}

func (h *InventoryHandler) ListLocationInventory(ctx context.Context, req *inventoryv1.ListLocationInventoryRequest) (*inventoryv1.ListLocationInventoryResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ListLocationInventory not implemented")
}

func (h *InventoryHandler) UpdateInventory(ctx context.Context, req *inventoryv1.UpdateInventoryRequest) (*inventoryv1.UpdateInventoryResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateInventory not implemented")
}
