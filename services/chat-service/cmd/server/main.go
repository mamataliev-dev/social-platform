package main

import (
	"context"
	"fmt"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/middleware"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	chatpb "github.com/mamataliev-dev/social-platform/api/gen/chat/v1"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/config"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/logger"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/mapper"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/repository"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/service"
)

// main is the entrypoint for the user-service application. It delegates startup
// logic to run and logs any fatal errors.
func main() {
	if err := run(); err != nil {
		log.Fatalf("application error: %v", err)
	}
}

// run loads configuration, sets up logging, initializes the database and all
// service dependencies, and starts both gRPC and HTTP REST servers. It handles
// graceful shutdown on interrupt signals. All dependencies are injected via
// constructors, following Dependency Inversion and Single Responsibility.
func run() error {
	// Load environment and configuration
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found, falling back to environment variables")
	}

	cfg, err := config.Load("config.yaml")
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		return err
	}

	// Initialize logger
	baseLogger := logger.SetupLogger(cfg.Env)
	slog.Info("logger initialized", "env", cfg.Env)

	// Database connection
	db, err := repository.NewPostgresConnection(cfg)
	if err != nil {
		slog.Error("failed to connect to Postgres", "error", err)
		return err
	}
	defer db.Close()

	// Initialize dependencies
	roomLogger := baseLogger.With("service", "room")

	mappers := mapper.NewMappers()

	roomRepo := repository.NewRoomPostgres(db, mappers.Room)

	roomSvc := service.NewRoomService(roomRepo, mappers.Room, roomLogger)

	// Context and signal handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// gRPC server setup
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.UnaryAuthInterceptor,
		),
	)

	chatpb.RegisterChatServiceServer(grpcServer, roomSvc)
	reflection.Register(grpcServer)

	grpcAddr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		slog.Error("failed to listen for gRPC", "error", err)
		return err
	}

	// HTTP REST Gateway setup
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := chatpb.RegisterChatServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		slog.Error("failed to register room gateway", "error", err)
		return err
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: mux,
	}

	// Run servers concurrently
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		slog.Info("starting gRPC server", "addr", grpcAddr)
		return grpcServer.Serve(lis)
	})
	eg.Go(func() error {
		slog.Info("starting HTTP REST server", "addr", httpServer.Addr)
		return httpServer.ListenAndServe()
	})

	// Wait for shutdown signal
	<-egCtx.Done()
	slog.Info("shutdown signal received, stopping servers...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcServer.GracefulStop()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown failed", "error", err)
	}

	// Wait for all goroutines to finish
	return eg.Wait()
}
