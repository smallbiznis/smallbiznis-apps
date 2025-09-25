package persistence

import (
	"context"

	"github.com/smallbiznis/smallbiznis-apps/internal/ledger/domain"
	"github.com/smallbiznis/smallbiznis-apps/pkg/db/option"
	"github.com/smallbiznis/smallbiznis-apps/pkg/repository"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type BalanceParams struct {
	fx.In
	DB *gorm.DB
}

type balanceRepository struct {
	db   *gorm.DB
	repo repository.Repository[domain.Balance]
}

func NewBalanceRepository(p BalanceParams) domain.BalanceRepository {
	return &balanceRepository{
		db:   p.DB,
		repo: repository.ProvideStore[domain.Balance](p.DB),
	}
}

func (r *balanceRepository) WithTrx(tx *gorm.DB) domain.BalanceRepository {
	return &balanceRepository{
		db:   tx,
		repo: repository.ProvideStore[domain.Balance](tx),
	}
}

func (r *balanceRepository) Count(ctx context.Context, f *domain.Balance) (int64, error) {
	return r.repo.Count(ctx, f)
}

func (r *balanceRepository) Find(ctx context.Context, f *domain.Balance, opts ...option.QueryOption) ([]*domain.Balance, error) {
	return r.repo.Find(ctx, f, opts...)
}

func (r *balanceRepository) FindOne(ctx context.Context, f *domain.Balance, opts ...option.QueryOption) (*domain.Balance, error) {
	return r.repo.FindOne(ctx, f, opts...)
}

func (r *balanceRepository) Create(ctx context.Context, entry *domain.Balance) error {
	return r.repo.Create(ctx, entry)
}

func (r *balanceRepository) Update(ctx context.Context, entryID string, entry any) error {
	return r.repo.Update(ctx, entryID, entry)
}
