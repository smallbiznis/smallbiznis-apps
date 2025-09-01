package organization

import (
	"github.com/smallbiznis/smallbiznis-apps/internal/organization/domain"
	"github.com/smallbiznis/smallbiznis-apps/internal/organization/usecase"
	"github.com/smallbiznis/smallbiznis-apps/pkg/repository"
	"go.uber.org/fx"
)

var Service = fx.Module("organization.service",
	fx.Provide(repository.ProvideStore[domain.Organization]), // Provide Repository
	fx.Provide(usecase.NewOrganizationUsecase),               // Provide Usecase/Business Logic
)
