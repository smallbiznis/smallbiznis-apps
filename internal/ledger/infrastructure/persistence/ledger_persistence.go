package persistence

import (
	"context"

	"github.com/smallbiznis/smallbiznis-apps/internal/ledger/domain"
	"github.com/smallbiznis/smallbiznis-apps/pkg/db/option"
	"github.com/smallbiznis/smallbiznis-apps/pkg/repository"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type LedgerParams struct {
	fx.In
	DB *gorm.DB
}

type ledgerRepository struct {
	db   *gorm.DB
	repo repository.Repository[domain.LedgerEntry]
}

func NewLedgerRepository(p LedgerParams) domain.LedgerRepository {
	return &ledgerRepository{
		db:   p.DB,
		repo: repository.ProvideStore[domain.LedgerEntry](p.DB),
	}
}

func (r *ledgerRepository) WithTrx(tx *gorm.DB) domain.LedgerRepository {
	return &ledgerRepository{
		db:   tx,
		repo: repository.ProvideStore[domain.LedgerEntry](tx),
	}
}

func (r *ledgerRepository) Count(ctx context.Context, f *domain.LedgerEntry) (int64, error) {
	return r.repo.Count(ctx, f)
}

func (r *ledgerRepository) Find(ctx context.Context, f *domain.LedgerEntry, opts ...option.QueryOption) ([]*domain.LedgerEntry, error) {
	return r.repo.Find(ctx, f, opts...)
}

func (r *ledgerRepository) FindOne(ctx context.Context, f *domain.LedgerEntry, opts ...option.QueryOption) (*domain.LedgerEntry, error) {
	return r.repo.FindOne(ctx, f, opts...)
}

func (r *ledgerRepository) Create(ctx context.Context, entry *domain.LedgerEntry) error {
	return r.repo.Create(ctx, entry)
}

func (r *ledgerRepository) Update(ctx context.Context, entryID string, entry *domain.LedgerEntry) error {
	return r.repo.Update(ctx, entryID, entry)
}
