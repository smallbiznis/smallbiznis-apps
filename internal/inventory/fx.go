package inventory

import (
	inventoryv1 "github.com/smallbiznis/go-genproto/smallbiznis/inventory"
	grpc_handler "github.com/smallbiznis/smallbiznis-apps/internal/inventory/interfaces/grpc"
	"github.com/smallbiznis/smallbiznis-apps/internal/inventory/usecase"
	"github.com/smallbiznis/smallbiznis-apps/pkg/server"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func RegisterServiceServer(s *grpc.Server, srv *grpc_handler.InventoryHandler) {
	inventoryv1.RegisterInventoryServiceServer(s, srv)
}

var Server = fx.Module("inventory.service",
	fx.Provide(
		server.NewListener,
		server.WithOption,
		server.NewGRPCServer,
		server.NewServeMux,
	),
	fx.Provide(
		usecase.NewInventory,
		grpc_handler.NewInventoryHandler,
	),
	fx.Invoke(
		RegisterServiceServer,
		server.StartGRPCServer,
	),
	server.NewServer,
)
