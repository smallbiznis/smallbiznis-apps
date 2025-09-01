package persistence

import (
	"context"

	"github.com/smallbiznis/smallbiznis-apps/internal/ledger/domain"
	"github.com/smallbiznis/smallbiznis-apps/pkg/db/option"
	"github.com/smallbiznis/smallbiznis-apps/pkg/repository"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type CreditPoolParams struct {
	fx.In
	DB *gorm.DB
}

type creditPoolRepository struct {
	db   *gorm.DB
	repo repository.Repository[domain.CreditPool]
}

func NewCreditPoolRepository(p CreditPoolParams) domain.CreditPoolRepository {
	return &creditPoolRepository{
		db:   p.DB,
		repo: repository.ProvideStore[domain.CreditPool](p.DB),
	}
}

func (r *creditPoolRepository) WithTrx(tx *gorm.DB) domain.CreditPoolRepository {
	return &creditPoolRepository{
		db:   tx,
		repo: repository.ProvideStore[domain.CreditPool](tx),
	}
}

func (r *creditPoolRepository) Count(ctx context.Context, f *domain.CreditPool) (int64, error) {
	return r.repo.Count(ctx, f)
}

func (r *creditPoolRepository) Find(ctx context.Context, f *domain.CreditPool, opts ...option.QueryOption) ([]*domain.CreditPool, error) {
	return r.repo.Find(ctx, f, opts...)
}

func (r *creditPoolRepository) FindOne(ctx context.Context, f *domain.CreditPool, opts ...option.QueryOption) (*domain.CreditPool, error) {
	return r.repo.FindOne(ctx, f, opts...)
}

func (r *creditPoolRepository) Create(ctx context.Context, entry *domain.CreditPool) error {
	return r.repo.Create(ctx, entry)
}

func (r *creditPoolRepository) Update(ctx context.Context, entryID string, entry any) error {
	return r.repo.Update(ctx, entryID, entry)
}
