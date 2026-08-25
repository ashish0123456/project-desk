package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	issuev1 "gitlab.com/Ash3941/project-management/gen/proto/issue/v1"
	"gitlab.com/Ash3941/project-management/pkg/logger"
	"gitlab.com/Ash3941/project-management/pkg/migrate"
	"gitlab.com/Ash3941/project-management/pkg/postgres"
	"gitlab.com/Ash3941/project-management/services/issue-service/internal/config"
	"gitlab.com/Ash3941/project-management/services/issue-service/internal/handler"
	"gitlab.com/Ash3941/project-management/services/issue-service/internal/server"
	"gitlab.com/Ash3941/project-management/services/issue-service/internal/service"
	"gitlab.com/Ash3941/project-management/services/issue-service/internal/store"
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
		log.Error("issue-service failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// apply migrations for issue service
	if err := migrate.ApplyMigrations(db, "services/issue-service/migrations"); err != nil {
		log.Error("issue-service failed to apply migrations", "error", err)
		os.Exit(1)
	}

	st := store.New(db)
	svc := service.New(st)
	h := handler.New(svc)

	httpHandler := h.Routes()
	srv := server.New(httpHandler, getEnv("ISSUE_HTTP_PORT", "8083"))

	go func() {
		log.Info("issue-service starting", "port", getEnv("ISSUE_HTTP_PORT", "8083"), "env", cfg.AppEnv)
		if err := srv.Start(); err != nil {
			log.Error("issue-service stopped", "error", err)
			os.Exit(1)
		}
	}()

	// start gRPC server for internal RPCs
	grpcServer := grpc.NewServer()
	issuev1.RegisterIssueServiceServer(grpcServer, &issueGRPCServer{svc: svc})
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Error("issue-service failed to listen on gRPC port", "error", err)
		os.Exit(1)
	}
	go func() {
		log.Info("issue-service grpc starting", "port", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("grpc server stopped", "error", err)
			os.Exit(1)
		}
	}()

	log.Info("issue-service started",
		"grpc_port", cfg.GRPCPort,
		"database", cfg.DB.DBName,
		"env", cfg.AppEnv,
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("issue-service shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err == nil {
		log.Info("issue-service database connection closed")
	}

	_ = srv.Shutdown(ctx)

	log.Info("issue-service stopped gracefully")
	// gracefully stop gRPC server
	grpcServer.GracefulStop()
}

type issueGRPCServer struct {
	issuev1.UnimplementedIssueServiceServer
	svc *service.Service
}

func (s *issueGRPCServer) CreateIssue(ctx context.Context, req *issuev1.CreateIssueRequest) (*issuev1.CreateIssueResponse, error) {
	issue, err := s.svc.CreateIssue(ctx, req.ProjectId, req.Title, req.Body)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create issue")
	}
	return &issuev1.CreateIssueResponse{
		Id:        issue.ID,
		ProjectId: issue.ProjectID,
		Title:     issue.Title,
		Body:      issue.Body,
		Status:    issue.Status,
		Priority:  issue.Priority,
		Assignee:  issue.Assignee,
		CreatedAt: issue.CreatedAt.Format(time.RFC3339),
	}, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
