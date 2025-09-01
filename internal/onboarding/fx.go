package ledger

import (
	ledgerv1 "github.com/smallbiznis/go-genproto/smallbiznis/ledger/v1"
	"github.com/smallbiznis/smallbiznis-apps/internal/ledger/infrastructure/persistence"
	grpc_handler "github.com/smallbiznis/smallbiznis-apps/internal/ledger/interfaces/grpc"
	"github.com/smallbiznis/smallbiznis-apps/internal/ledger/usecase"
	"github.com/smallbiznis/smallbiznis-apps/pkg/server"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func RegisterServiceServer(s *grpc.Server, srv *grpc_handler.Handler) {
	ledgerv1.RegisterLedgerServiceServer(s, srv)
}

// func RegisterServiceHandlerFromEndpoint(lc fx.Lifecycle, mux *runtime.ServeMux, cfg *config.Config) {
// 	lc.Append(fx.Hook{
// 		OnStart: func(ctx context.Context) error {

// 			opts := []grpc.DialOption{
// 				grpc.WithTransportCredentials(insecure.NewCredentials()),
// 			}

// 			if err := ledgerv1.RegisterServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%s", cfg.Grpc.Addr), opts); err != nil {
// 				zap.L().Error("failed to RegisterServiceHandlerFromEndpoint", zap.Error(err))
// 			}

// 			return nil
// 		},
// 	})
// }

var Server = fx.Module("rulengine.service.server",
	fx.Provide(
		server.NewListener,
		server.WithOption,
		server.NewGRPCServer,
		server.NewServeMux,
	),
	fx.Provide(
		persistence.NewLedgerRepository,
		persistence.NewCreditPoolRepository,
		persistence.NewBalanceRepository,
		usecase.NewLedger,
		grpc_handler.NewHandler,
	),
	fx.Invoke(
		RegisterServiceServer,
		// RegisterServiceHandlerFromEndpoint,
		server.StartGRPCServer,
	),
	server.NewServer,
)
