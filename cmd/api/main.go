package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/Ash3941/go-pm-platform/internal/config"
	"gitlab.com/Ash3941/go-pm-platform/internal/database"
	"gitlab.com/Ash3941/go-pm-platform/internal/handler"
	"gitlab.com/Ash3941/go-pm-platform/internal/logger"
	"gitlab.com/Ash3941/go-pm-platform/internal/router"
	"gitlab.com/Ash3941/go-pm-platform/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.LogLevel)

	db, err := database.New(cfg)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	defer db.Close()

	systemHandler := handler.NewSystemHandler(cfg, db)
	router := router.New(systemHandler, log)
	server := server.New(router, cfg.AppPort)

	go func() {
		log.Info("starting server", "port", cfg.AppPort)
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server stopped", "error", err)
			os.Exit(1)
		} 
	}()
	
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	
	<-stop

	log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed", "error", err)
		return
	}

	log.Info("server stopped gracefully")
}