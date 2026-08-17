package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string
	LogLevel string
	Version string

	DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string
}

func LoadConfig() (*Config, error) {

	_ = godotenv.Load()

	cfg := & Config{
		AppName: getEnv("APP_NAME", "Project Management Platform"),
		AppEnv: getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Version: getEnv("VERSION", "1.0.0"),

		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName: getEnv("DB_NAME", "project_management"),
		DBSSLMode: getEnv("DB_SSLMODE", "disable"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}