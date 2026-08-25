package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"gitlab.com/Ash3941/project-management/pkg/logger"
	"gitlab.com/Ash3941/project-management/services/worker-service/internal/config"
)

func main() {
	// Worker process: simple Redis-backed worker pool that BLPOP's on the
	// "jobs" list and prints job payloads. In production you'd dispatch
	// to job handlers, implement retries, observability, and idempotency.
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.LogLevel)

	log.Info("worker-service started",
		"workers", cfg.WorkerCount,
		"redis", cfg.RedisAddr,
		"env", cfg.AppEnv,
	)
	// connect to redis
	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})

	workerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start workers
	for i := 0; i < cfg.WorkerCount; i++ {
		go func(id int) {
			for {
				res, err := rdb.BLPop(workerCtx, 0*time.Second, "jobs").Result()
				if err != nil {
					if err == context.Canceled || err == redis.Nil {
						return
					}
					fmt.Printf("worker %d: redis error: %v", id, err)
					time.Sleep(1 * time.Second)
					continue
				}
				if len(res) >= 2 {
					job := res[1]
					fmt.Printf("worker %d processing job: %s\n", id, job)
				}
			}
		}(i)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("worker-service shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = shutdownCtx

	log.Info("worker-service stopped gracefully")
}
