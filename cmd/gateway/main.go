// Package main is the entry point for the AI Gateway server.
// It initializes configuration, sets up the Envoy external processor,
// and starts the gRPC server to handle AI provider routing.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	grpcAddr = flag.String("grpc-addr", ":8080", "Address for the gRPC external processor server")
	logLevel = flag.String("log-level", "info", "Log level: debug, info, warn, error")
)

func main() {
	flag.Parse()

	logger := newLogger(*logLevel)
	slog.SetDefault(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx); err != nil {
		slog.Error("gateway exited with error", "error", err)
		os.Exit(1)
	}
	slog.Info("gateway shutdown complete")
}

func run(ctx context.Context) error {
	lis, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", *grpcAddr, err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(loggingUnaryInterceptor),
		grpc.ChainStreamInterceptor(loggingStreamInterceptor),
	)

	// Register health service
	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)
	healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Enable gRPC reflection for debugging
	reflection.Register(grpcServer)

	slog.Info("starting AI gateway gRPC server", "addr", *grpcAddr)

	errCh := make(chan error, 1)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			errCh <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutting down gRPC server")
		grpcServer.GracefulStop()
		return nil
	case err := <-errCh:
		return err
	}
}

func newLogger(level string) *slog.Logger {
	var l slog.Level
	switch level {
	case "debug":
		l = slog.LevelDebug
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l}))
}

func loggingUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	slog.Debug("unary RPC", "method", info.FullMethod)
	return handler(ctx, req)
}

func loggingStreamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	slog.Debug("stream RPC", "method", info.FullMethod)
	return handler(srv, ss)
}
