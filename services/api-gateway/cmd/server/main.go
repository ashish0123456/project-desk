package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/Ash3941/project-management/pkg/logger"
	"gitlab.com/Ash3941/project-management/services/api-gateway/internal/config"
	"gitlab.com/Ash3941/project-management/services/api-gateway/internal/handler"
	"gitlab.com/Ash3941/project-management/services/api-gateway/internal/router"
	"gitlab.com/Ash3941/project-management/services/api-gateway/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.LogLevel)

	systemHandler := handler.NewSystemHandler(cfg)
	httpHandler := router.New(systemHandler, log)
	srv := server.New(httpHandler, cfg.AppPort)

	go func() {
		log.Info("api-gateway starting", "port", cfg.AppPort, "env", cfg.AppEnv)
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("api-gateway stopped", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("api-gateway shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("api-gateway shutdown failed", "error", err)
		os.Exit(1)
	}

	log.Info("api-gateway stopped gracefully")
}
