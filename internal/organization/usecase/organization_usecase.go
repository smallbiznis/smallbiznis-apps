package usecase

import (
	"context"

	orgv1 "github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/smallbiznis-apps/internal/organization/domain"
	"github.com/smallbiznis/smallbiznis-apps/pkg/repository"
	"gorm.io/gorm"
)

type IOrganizationUsecase interface {
	CreateOrganization(context.Context, *orgv1.CreateOrganizationRequest) (*orgv1.Organization, error)
	GetOrganization(context.Context, *orgv1.GetOrganizationRequest) (*orgv1.Organization, error)
	ListOrganization(context.Context, *orgv1.ListOrganizationRequest) (*orgv1.ListOrganizationResponse, error)
	UpdateOrganization(context.Context, *orgv1.UpdateOrganizationRequest) (*orgv1.Organization, error)
}

type organizationUsecase struct {
	db           *gorm.DB
	countryRepo  repository.Repository[domain.Country]
	timezoneRepo repository.Repository[domain.Timezone]
	currencyRepo repository.Repository[domain.Currency]
	orgRepo      repository.Repository[domain.Organization]
}

func NewOrganizationUsecase(
	db *gorm.DB,
	countryRepo repository.Repository[domain.Country],
	timezoneRepo repository.Repository[domain.Timezone],
	currencyRepo repository.Repository[domain.Currency],
	orgRepo repository.Repository[domain.Organization],
) IOrganizationUsecase {
	return &organizationUsecase{
		db:           db,
		countryRepo:  countryRepo,
		timezoneRepo: timezoneRepo,
		currencyRepo: currencyRepo,
		orgRepo:      orgRepo,
	}
}

func (uc *organizationUsecase) LookupOrganization(ctx context.Context, req *orgv1.CreateOrganizationRequest) (*orgv1.Organization, error) {
}

func (uc *organizationUsecase) CreateOrganization(ctx context.Context, req *orgv1.CreateOrganizationRequest) (*orgv1.Organization, error) {
}

func (uc *organizationUsecase) GetOrganization(ctx context.Context, req *orgv1.GetOrganizationRequest) (*orgv1.Organization, error) {
}

func (uc *organizationUsecase) ListOrganization(ctx context.Context, req *orgv1.ListOrganizationRequest) (*orgv1.ListOrganizationResponse, error) {
}

func (uc *organizationUsecase) UpdateOrganization(ctx context.Context, req *orgv1.UpdateOrganizationRequest) (*orgv1.Organization, error) {
}
