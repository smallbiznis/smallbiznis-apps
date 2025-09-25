package usecase

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
)

//go:generate mockgen -source=inventory_usecase.go -destination=./../../usecase/mock_inventory_usecase.go -package=usecase
type InventoryUsecase interface {
}

type InventoryParams struct {
	fx.In
	DB *gorm.DB
}

type inventoryUsecase struct {
	db *gorm.DB
}

func NewInventory(p InventoryParams) InventoryUsecase {
	return &inventoryUsecase{
		db: p.DB,
	}
}
