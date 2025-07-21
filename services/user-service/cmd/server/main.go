// Package main is the entrypoint for the user-service. It loads configuration,
// initializes dependencies (logger, database, repositories, services), and starts
// both gRPC and HTTP REST servers. This package follows the Dependency Inversion
// principle by injecting abstractions for repositories, services, and utilities.
// It is closed for modification but open for extension via configuration and DI.
package main

import (
	"context"
	"fmt"
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

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/config"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/logger"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/mapper"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/middleware"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/repository"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/security"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/service"
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
	logger.SetupLogger(cfg.Env)
	slog.Info("logger initialized", "env", cfg.Env)

	// Database connection
	db, err := repository.NewPostgresConnection(cfg)
	if err != nil {
		slog.Error("failed to connect to Postgres", "error", err)
		return err
	}
	defer db.Close()

	// Initialize dependencies
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	tokenTTL := 15 * time.Minute

	hasher := security.BcryptHasher{}
	jwtGen := security.NewJWTGenerator(secretKey, tokenTTL)
	converter := mapper.NewMapper()

	userRepo := repository.NewUserPostgres(db)
	authRepo := repository.NewAuthPostgres(db)
	tokenRepo := repository.NewTokenPostgres(db)

	authSvc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtGen, hasher, converter)
	publicUserSvc := service.NewUserService(userRepo, converter)
	internalUserSvc := service.NewInternalUserService(userRepo, converter)

	// Context and signal handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// gRPC server setup
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.ValidationInterceptor(),
			middleware.TimeoutInterceptor,
			middleware.UnaryAuthInterceptor,
		),
	)

	userauthpb.RegisterAuthServiceServer(grpcServer, authSvc)
	userpb.RegisterUserServiceServer(grpcServer, publicUserSvc)
	userpb.RegisterInternalUserServiceServer(grpcServer, internalUserSvc)
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
	if err := userauthpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		slog.Error("failed to register auth gateway", "error", err)
		return err
	}
	if err := userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		slog.Error("failed to register user gateway", "error", err)
		return err
	}
	if err := userpb.RegisterInternalUserServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		slog.Error("failed to register internal user gateway", "error", err)
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
