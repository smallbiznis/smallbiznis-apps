package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/smallbiznis/smallbiznis-apps/pkg/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var ProvideGRPCServer = fx.Module("grpc.server",
	fx.Provide(NewListener),
	fx.Invoke(StartGRPCServer),
)

func NewListener(cfg *config.Config) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf(":%s", cfg.Grpc.Addr))
}

func WithOption(opts ...grpc.ServerOption) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			validator.UnaryServerInterceptor(validator.WithFailFast()),
		),
		grpc.ChainStreamInterceptor(
			validator.StreamServerInterceptor(validator.WithFailFast()),
		),
	}
}

func WithStatsHandler(tp trace.TracerProvider, mp metric.MeterProvider) grpc.ServerOption {
	return grpc.StatsHandler(
		otelgrpc.NewServerHandler(
			otelgrpc.WithTracerProvider(tp),
			otelgrpc.WithMeterProvider(mp),
		),
	)
}

// LoadCertificate
func LoadCertificate(certPath, keyPath string) (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

// WithTLS
func WithTLS(tls *tls.Certificate) grpc.ServerOption {
	return grpc.Creds(
		credentials.NewServerTLSFromCert(tls),
	)
}

func NewGRPCServer(opts ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(opts...)
}

func StartGRPCServer(lc fx.Lifecycle, lis net.Listener, srv *grpc.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				zap.L().Info("Starting gRPC server", zap.String("addr", lis.Addr().String()))
				if err := srv.Serve(lis); err != nil {
					zap.L().Fatal("gRPC server exited", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			zap.L().Info("Stopping gRPC server")
			srv.GracefulStop()
			return nil
		},
	})
}
