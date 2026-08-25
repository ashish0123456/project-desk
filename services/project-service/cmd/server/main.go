package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	projectv1 "gitlab.com/Ash3941/project-management/gen/proto/project/v1"
	"gitlab.com/Ash3941/project-management/pkg/logger"
	"gitlab.com/Ash3941/project-management/pkg/migrate"
	"gitlab.com/Ash3941/project-management/pkg/postgres"
	"gitlab.com/Ash3941/project-management/services/project-service/internal/config"
	"gitlab.com/Ash3941/project-management/services/project-service/internal/handler"
	"gitlab.com/Ash3941/project-management/services/project-service/internal/server"
	"gitlab.com/Ash3941/project-management/services/project-service/internal/service"
	"gitlab.com/Ash3941/project-management/services/project-service/internal/store"
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
		log.Error("project-service failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// apply migrations for project service
	if err := migrate.ApplyMigrations(db, "services/project-service/migrations"); err != nil {
		log.Error("project-service failed to apply migrations", "error", err)
		os.Exit(1)
	}

	st := store.New(db)
	svc := service.New(st)
	h := handler.New(svc)

	httpHandler := h.Routes()
	srv := server.New(httpHandler, getEnv("PROJECT_HTTP_PORT", "8082"))

	go func() {
		log.Info("project-service starting", "port", getEnv("PROJECT_HTTP_PORT", "8082"), "env", cfg.AppEnv)
		if err := srv.Start(); err != nil {
			log.Error("project-service stopped", "error", err)
			os.Exit(1)
		}
	}()

	// start gRPC server for internal RPCs
	grpcServer := grpc.NewServer()
	projectv1.RegisterProjectServiceServer(grpcServer, &projectGRPCServer{svc: svc})
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Error("project-service failed to listen on gRPC port", "error", err)
		os.Exit(1)
	}
	go func() {
		log.Info("project-service grpc starting", "port", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("grpc server stopped", "error", err)
			os.Exit(1)
		}
	}()

	log.Info("project-service started",
		"grpc_port", cfg.GRPCPort,
		"database", cfg.DB.DBName,
		"env", cfg.AppEnv,
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("project-service shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err == nil {
		log.Info("project-service database connection closed")
	}

	_ = srv.Shutdown(ctx)

	log.Info("project-service stopped gracefully")
	// gracefully stop gRPC server
	grpcServer.GracefulStop()
}

type projectGRPCServer struct {
	projectv1.UnimplementedProjectServiceServer
	svc *service.Service
}

func (s *projectGRPCServer) CreateProject(ctx context.Context, req *projectv1.CreateProjectRequest) (*projectv1.CreateProjectResponse, error) {
	project, err := s.svc.CreateProject(ctx, req.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create project")
	}
	return &projectv1.CreateProjectResponse{
		Id:        project.ID,
		Name:      project.Name,
		CreatedAt: project.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *projectGRPCServer) AddProjectMember(ctx context.Context, req *projectv1.AddProjectMemberRequest) (*projectv1.AddProjectMemberResponse, error) {
	member, err := s.svc.AddMember(ctx, req.ProjectId, req.UserId, req.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to add project member")
	}
	return &projectv1.AddProjectMemberResponse{
		Id:        member.ID,
		ProjectId: member.ProjectID,
		UserId:    member.UserID,
		Role:      member.Role,
		CreatedAt: member.CreatedAt.Format(time.RFC3339),
	}, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
