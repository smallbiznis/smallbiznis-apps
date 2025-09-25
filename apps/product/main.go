package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/smallbiznis/smallbiznis-apps/internal/product"
	"github.com/smallbiznis/smallbiznis-apps/pkg/config"
	"github.com/smallbiznis/smallbiznis-apps/pkg/db"
	"github.com/smallbiznis/smallbiznis-apps/pkg/logger"
	"github.com/smallbiznis/smallbiznis-apps/pkg/pprof"
	"github.com/smallbiznis/smallbiznis-apps/pkg/profiling"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var fxLogger = fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
	return &fxevent.ZapLogger{Logger: logger}
})

var zapField = fx.Provide(func(cfg *config.Config) []zap.Field {
	return []zap.Field{
		zap.String("app_name", cfg.AppName),
		zap.String("app_version", cfg.AppVersion),
		zap.String("app_env", cfg.AppEnv),
	}
})

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := fx.New(
		config.Module,
		zapField,
		db.Module,
		logger.Module,
		profiling.Module,
		product.Server,
		pprof.Module,
		fxLogger,
	)

	if err := app.Start(ctx); err != nil {
		zap.L().Fatal("Failed to start", zap.Error(err))
	}

	// Wait for signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	zap.L().Info("Server is shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := app.Stop(shutdownCtx); err != nil {
		zap.L().Error("Failed to stop", zap.Error(err))
	}
}
