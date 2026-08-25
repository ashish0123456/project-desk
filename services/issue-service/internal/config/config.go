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
		GRPCPort: getEnv("ISSUE_GRPC_PORT", "50053"),
		DB: postgres.Config{
			Host:     getEnv("ISSUE_DB_HOST", "localhost"),
			Port:     getEnv("ISSUE_DB_PORT", "5434"),
			User:     getEnv("ISSUE_DB_USER", "postgres"),
			Password: getEnv("ISSUE_DB_PASSWORD", "postgres"),
			DBName:   getEnv("ISSUE_DB_NAME", "pm_issue"),
			SSLMode:  getEnv("ISSUE_DB_SSLMODE", "disable"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
