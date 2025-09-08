package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct { // 1.3 need config data
	ServerPort  string
	DataBaseUrl string
	Environment string
	LogLevel    string
}

func LoadConfig() (*Config, error) { // 1.4 load env variable using godotenv
	if err := godotenv.Load(); err != nil {

		return nil, fmt.Errorf("Error loading file : %v ", err)

	}
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DataBaseUrl: getEnv("DATABASE_URL", "postgres"),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}, nil
}

func getEnv(key, defaultValue string) string { // 1.5 more safety if Load env failed
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
