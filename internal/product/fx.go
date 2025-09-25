package product

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	productv1 "github.com/smallbiznis/go-genproto/smallbiznis/product"
	grpc_handler "github.com/smallbiznis/smallbiznis-apps/internal/product/interfaces/grpc"
	"github.com/smallbiznis/smallbiznis-apps/internal/product/usecase"
	"github.com/smallbiznis/smallbiznis-apps/pkg/config"
	"github.com/smallbiznis/smallbiznis-apps/pkg/gen"
	"github.com/smallbiznis/smallbiznis-apps/pkg/server"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RegisterServiceServer(s *grpc.Server, srv *grpc_handler.ProductHandler) {
	productv1.RegisterProductServiceServer(s, srv)
}

func RegisterServiceHandlerFromEndpoint(lc fx.Lifecycle, mux *runtime.ServeMux, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			opts := []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			}

			if err := productv1.RegisterProductServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%s", cfg.Grpc.Addr), opts); err != nil {
				zap.L().Error("failed to RegisterServiceHandlerFromEndpoint", zap.Error(err))
			}

			return nil
		},
	})
}

var Server = fx.Module("product.service",
	fx.Provide(
		server.NewListener,
		server.WithOption,
		server.NewGRPCServer,
		server.NewServeMux,
	),
	fx.Provide(
		grpc_handler.NewProductHandler,
		gen.NewSnowflakeNode,
		usecase.NewProduct,
	),
	fx.Invoke(
		RegisterServiceServer,
		RegisterServiceHandlerFromEndpoint,
		server.StartGRPCServer,
	),
	server.NewServer,
)
