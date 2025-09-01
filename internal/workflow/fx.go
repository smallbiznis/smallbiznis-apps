package workflow

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	workflowv1 "github.com/smallbiznis/go-genproto/smallbiznis/workflow/v1"
	"github.com/smallbiznis/smallbiznis-apps/internal/workflow/domain"
	grpc_handler "github.com/smallbiznis/smallbiznis-apps/internal/workflow/interfaces/grpc"
	"github.com/smallbiznis/smallbiznis-apps/internal/workflow/usecase"
	"github.com/smallbiznis/smallbiznis-apps/pkg/config"
	"github.com/smallbiznis/smallbiznis-apps/pkg/repository"
	"github.com/smallbiznis/smallbiznis-apps/pkg/server"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RegisterServiceServer(s *grpc.Server, srv *grpc_handler.FlowHandler) {
	workflowv1.RegisterWorkflowServiceServer(s, srv)
}

func RegisterServiceHandlerFromEndpoint(lc fx.Lifecycle, mux *runtime.ServeMux, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			opts := []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			}

			if err := workflowv1.RegisterWorkflowServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%s", cfg.Grpc.Addr), opts); err != nil {
				zap.L().Error("failed to RegisterWorkflowServiceHandlerFromEndpoint", zap.Error(err))
			}

			return nil
		},
	})
}

var Server = fx.Module("workflow.service.server",
	fx.Provide(
		server.NewListener,
		server.WithOption,
		server.NewGRPCServer,
		server.NewServeMux,
	),
	fx.Provide(
		repository.ProvideStore[domain.FlowTemplate],
		repository.ProvideStore[domain.Flow],
		repository.ProvideStore[domain.Node],
		repository.ProvideStore[domain.Edge],
		usecase.NewFlowUsecase,
		grpc_handler.NewFlowHandler,
	),
	fx.Invoke(
		RegisterServiceServer,
		RegisterServiceHandlerFromEndpoint,
		server.StartGRPCServer,
	),
	server.NewServer,
)
