package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/config"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/logger"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/middleware"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/repository"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.SetupLogger(cfg.Env)
	slog.Info("Logger initialized", "env", cfg.Env)

	// Setup Postgres connection + repository
	db, err := repository.NewPostgresConnection(cfg)
	if err != nil {
		slog.Error("Failed to connect to Postgres", "error", err)
		os.Exit(1)
	}

	// Initialize repository + service
	authRepo := repository.NewAuthPostgres(db)
	authService := service.NewAuthService(authRepo)

	// Start gRPC server
	// TODO: Connect TLS
	grpcAddr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		slog.Error("Failed to listen for gRPC", "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.TimeoutInterceptor))
	userauthpb.RegisterAuthServiceServer(grpcServer, authService)

	go func() {
		slog.Info("Starting gRPC server", "address", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server failed", "error", err)
			os.Exit(1)
		}
	}()

	// REST Gateway
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	_ = userauthpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)

	log.Println("HTTP REST on", ":100")
	if err := http.ListenAndServe(":100", mux); err != nil {
		slog.Error("HTTP REST failed", "error", err)
		os.Exit(1)
	}
}
