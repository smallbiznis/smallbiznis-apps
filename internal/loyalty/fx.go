package point

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pointv1 "github.com/smallbiznis/go-genproto/smallbiznis/loyalty/point/v1"
	"github.com/smallbiznis/smallbiznis-apps/internal/loyalty/domain"
	grpc_handler "github.com/smallbiznis/smallbiznis-apps/internal/loyalty/interfaces/grpc"
	"github.com/smallbiznis/smallbiznis-apps/internal/loyalty/usecase"
	"github.com/smallbiznis/smallbiznis-apps/pkg/config"
	"github.com/smallbiznis/smallbiznis-apps/pkg/repository"
	"github.com/smallbiznis/smallbiznis-apps/pkg/server"
	"github.com/smallbiznis/smallbiznis-apps/pkg/workflow"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RegisterServiceServer(s *grpc.Server, srv *grpc_handler.PointHandler) {
	pointv1.RegisterPointServiceServer(s, srv)
}

func RegisterServiceHandlerFromEndpoint(lc fx.Lifecycle, mux *runtime.ServeMux, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			opts := []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			}
			if err := pointv1.RegisterPointServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%s", cfg.Grpc.Addr), opts); err != nil {
				zap.L().Error("failed to RegisterPointServiceHandlerFromEndpoint", zap.Error(err))
				return err
			}

			zap.L().Info("Success RegisterServiceHandlerFromEndoint")
			return nil
		},
	})
}

var Service = fx.Module("point.service",
	fx.Provide(
		server.NewListener,
		server.WithOption,
		server.NewGRPCServer,
		server.NewServeMux,
	),
	workflow.ProvideClient,
	fx.Provide(
		repository.ProvideStore[domain.Transaction],
		usecase.NewPointUsecase,
		grpc_handler.NewPointHandler,
	),
	fx.Invoke(
		RegisterServiceServer,
		RegisterServiceHandlerFromEndpoint,
		server.StartGRPCServer,
	),
	server.NewServer,
)
