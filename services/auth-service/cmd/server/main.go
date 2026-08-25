package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authv1 "gitlab.com/Ash3941/project-management/gen/proto/auth/v1"
	"gitlab.com/Ash3941/project-management/pkg/logger"
	"gitlab.com/Ash3941/project-management/pkg/migrate"
	"gitlab.com/Ash3941/project-management/pkg/postgres"
	"gitlab.com/Ash3941/project-management/services/auth-service/internal/config"
	"gitlab.com/Ash3941/project-management/services/auth-service/internal/handler"
	"gitlab.com/Ash3941/project-management/services/auth-service/internal/jwt"
	"gitlab.com/Ash3941/project-management/services/auth-service/internal/server"
	"gitlab.com/Ash3941/project-management/services/auth-service/internal/service"
	"gitlab.com/Ash3941/project-management/services/auth-service/internal/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.LogLevel)

	db, err := postgres.Open(cfg.DB)
	if err != nil {
		log.Error("auth-service failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// apply SQL migrations for auth service
	if err := migrate.ApplyMigrations(db, "services/auth-service/migrations"); err != nil {
		log.Error("auth-service failed to apply migrations", "error", err)
		os.Exit(1)
	}

	st := store.New(db)
	svc := service.New(st)
	jwtm := jwt.New(getEnv("AUTH_JWT_SECRET", "devsecret"))
	h := handler.New(svc, jwtm)

	httpHandler := h.Routes()
	srv := server.New(httpHandler, getEnv("AUTH_HTTP_PORT", "8081"))

	// start gRPC server for internal service-to-service RPCs
	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, &authGRPCServer{svc: svc, jwtm: jwtm})
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Error("auth-service failed to listen on gRPC port", "error", err)
		os.Exit(1)
	}
	go func() {
		log.Info("auth-service grpc starting", "port", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("grpc server stopped", "error", err)
			os.Exit(1)
		}
	}()

	go func() {
		log.Info("auth-service starting", "port", getEnv("AUTH_HTTP_PORT", "8081"), "env", cfg.AppEnv)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Error("auth-service stopped", "error", err)
			os.Exit(1)
		}
	}()

	log.Info("auth-service started",
		"grpc_port", cfg.GRPCPort,
		"database", cfg.DB.DBName,
		"env", cfg.AppEnv,
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("auth-service shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err == nil {
		log.Info("auth-service database connection closed")
	}

	_ = srv.Shutdown(ctx)

	// gracefully stop gRPC server
	grpcServer.GracefulStop()

	log.Info("auth-service stopped gracefully")
}

type authGRPCServer struct {
	authv1.UnimplementedAuthServiceServer
	svc  *service.Service
	jwtm *jwt.Manager
}

func (s *authGRPCServer) Register(ctx context.Context, req *authv1.RegisterUserRequest) (*authv1.RegisterUserResponse, error) {
	user, err := s.svc.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to register user")
	}
	return &authv1.RegisterUserResponse{Id: user.ID, Email: user.Email}, nil
}

func (s *authGRPCServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	user, err := s.svc.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	access, err := s.jwtm.NewAccessToken(user.ID, 15*time.Minute)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create access token")
	}
	refresh, err := s.jwtm.NewRefreshToken(user.ID, 7*24*time.Hour)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create refresh token")
	}
	return &authv1.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    time.Now().Add(15 * time.Minute).UTC().Format(time.RFC3339),
	}, nil
}

func (s *authGRPCServer) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	claims, err := s.jwtm.Validate(req.Token)
	if err != nil {
		return &authv1.ValidateTokenResponse{Valid: false, Error: err.Error()}, nil
	}
	return &authv1.ValidateTokenResponse{Valid: true, UserId: claims.UserID}, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
