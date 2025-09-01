package domain

import (
	"context"

	"github.com/smallbiznis/smallbiznis-apps/pkg/db/option"
	"gorm.io/gorm"
)

//go:generate mockgen -source=repository.go -destination=./../../repository/mock_ledger_repository.go -package=repository
type LedgerRepository interface {
	WithTrx(tx *gorm.DB) LedgerRepository
	Find(ctx context.Context, query *LedgerEntry, opts ...option.QueryOption) ([]*LedgerEntry, error)
	FindOne(ctx context.Context, query *LedgerEntry, opts ...option.QueryOption) (*LedgerEntry, error)
	Create(ctx context.Context, resource *LedgerEntry) error
	// Update(ctx context.Context, resourceID string, resource *LedgerEntry) error
	// Delete(ctx context.Context, resourceID string) error
	// BatchCreate(ctx context.Context, resources []*LedgerEntry, batchSize int) error
	// BatchUpdate(ctx context.Context, resources []*LedgerEntry) error
	Count(ctx context.Context, query *LedgerEntry) (int64, error)
}

type CreditPoolRepository interface {
	WithTrx(tx *gorm.DB) CreditPoolRepository
	Find(ctx context.Context, query *CreditPool, opts ...option.QueryOption) ([]*CreditPool, error)
	FindOne(ctx context.Context, query *CreditPool, opts ...option.QueryOption) (*CreditPool, error)
	Create(ctx context.Context, resource *CreditPool) error
	Update(ctx context.Context, resourceID string, resource any) error
	// Delete(ctx context.Context, resourceID string) error
	// BatchCreate(ctx context.Context, resources []*CreditPool, batchSize int) error
	// BatchUpdate(ctx context.Context, resources []*CreditPool) error
	Count(ctx context.Context, query *CreditPool) (int64, error)
}

type BalanceRepository interface {
	WithTrx(tx *gorm.DB) BalanceRepository
	Find(ctx context.Context, query *Balance, opts ...option.QueryOption) ([]*Balance, error)
	FindOne(ctx context.Context, query *Balance, opts ...option.QueryOption) (*Balance, error)
	Create(ctx context.Context, resource *Balance) error
	Update(ctx context.Context, resourceID string, resource any) error
	// Delete(ctx context.Context, resourceID string) error
	// BatchCreate(ctx context.Context, resources []*Balance, batchSize int) error
	// BatchUpdate(ctx context.Context, resources []*Balance) error
	Count(ctx context.Context, query *Balance) (int64, error)
}
