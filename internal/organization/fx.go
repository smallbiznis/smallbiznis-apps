package organization

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	orgv1 "github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	grpc_handler "github.com/smallbiznis/smallbiznis-apps/internal/organization/interfaces/grpc"
	"github.com/smallbiznis/smallbiznis-apps/internal/organization/usecase"
	"github.com/smallbiznis/smallbiznis-apps/pkg/config"
	"github.com/smallbiznis/smallbiznis-apps/pkg/gen"
	"github.com/smallbiznis/smallbiznis-apps/pkg/server"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RegisterServiceServer(s *grpc.Server, srv *grpc_handler.OrganizationHandler) {
	orgv1.RegisterOrganizationServiceServer(s, srv)
}

func RegisterServiceHandlerFromEndpoint(lc fx.Lifecycle, mux *runtime.ServeMux, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			opts := []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			}

			if err := orgv1.RegisterOrganizationServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%s", cfg.Grpc.Addr), opts); err != nil {
				zap.L().Error("failed to RegisterWorkflowServiceHandlerFromEndpoint", zap.Error(err))
			}

			return nil
		},
	})
}

var Server = fx.Module("organization.service",
	fx.Provide(
		server.NewListener,
		server.WithOption,
		server.NewGRPCServer,
		server.NewServeMux,
	),
	fx.Provide(
		usecase.NewCountry,
		usecase.NewOrganization,
	),
	fx.Provide(
		grpc_handler.NewOrganization,
		gen.NewSnowflakeNode,
	),
	fx.Invoke(
		RegisterServiceServer,
		RegisterServiceHandlerFromEndpoint,
		server.StartGRPCServer,
	),
	server.NewServer,
)
