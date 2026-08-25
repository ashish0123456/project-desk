package config

import (
	"os"

	"github.com/joho/godotenv"
	"gitlab.com/Ash3941/project-management/pkg/postgres"
)

type Config struct {
	AppEnv   string
	LogLevel string
	GRPCPort string
	DB       postgres.Config
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		AppEnv:   getEnv("APP_ENV", "development"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		GRPCPort: getEnv("PROJECT_GRPC_PORT", "50052"),
		DB: postgres.Config{
			Host:     getEnv("PROJECT_DB_HOST", "localhost"),
			Port:     getEnv("PROJECT_DB_PORT", "5433"),
			User:     getEnv("PROJECT_DB_USER", "postgres"),
			Password: getEnv("PROJECT_DB_PASSWORD", "postgres"),
			DBName:   getEnv("PROJECT_DB_NAME", "pm_project"),
			SSLMode:  getEnv("PROJECT_DB_SSLMODE", "disable"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
