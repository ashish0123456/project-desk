package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName  string
	AppEnv   string
	AppPort  string
	LogLevel string
	Version  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		AppName:  getEnv("APP_NAME", "Project Management Platform"),
		AppEnv:   getEnv("APP_ENV", "development"),
		AppPort:  getEnv("GATEWAY_PORT", getEnv("APP_PORT", "8080")),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Version:  getEnv("VERSION", "0.1.0"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
