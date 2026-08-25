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
		GRPCPort: getEnv("AUTH_GRPC_PORT", "50051"),
		DB: postgres.Config{
			Host:     getEnv("AUTH_DB_HOST", "localhost"),
			Port:     getEnv("AUTH_DB_PORT", "5432"),
			User:     getEnv("AUTH_DB_USER", "postgres"),
			Password: getEnv("AUTH_DB_PASSWORD", "postgres"),
			DBName:   getEnv("AUTH_DB_NAME", "pm_auth"),
			SSLMode:  getEnv("AUTH_DB_SSLMODE", "disable"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
