package organization

import (
	"github.com/smallbiznis/smallbiznis-apps/internal/organization/usecase"
	"go.uber.org/fx"
)

var Service = fx.Module("organization.service",
	fx.Provide(usecase.NewOrganizationUsecase), // Provide Usecase/Business Logic
)
