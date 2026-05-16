package utils

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"cascade/pkg/logger"
)

type Config struct {
	ApplicationPort int
	ApplicationUrl  string
	AllowedOrigins  []string

	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logger.Error("No .env file found, using environment variables", nil)
	}

	applicationPort, err := strconv.Atoi(getEnv("APPLICATION_PORT", "4000"))
	if err != nil {
		return nil, err
	}

	postgresPort, err := strconv.Atoi(getEnv("POSTGRES_PORT", "5433"))
	if err != nil {
		return nil, err
	}

	allowedOrigins := getEnv("ALLOWED_ORIGINS", "")
	allowedOriginsList := strings.Split(allowedOrigins, ",")

	config := &Config{
		ApplicationPort: applicationPort,
		ApplicationUrl:  os.ExpandEnv(getEnv("APPLICATION_URL", "")),
		AllowedOrigins:  allowedOriginsList,

		PostgresUser:     getEnv("POSTGRES_USER", ""),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
		PostgresHost:     getEnv("POSTGRES_HOST", ""),
		PostgresPort:     postgresPort,
		PostgresDatabase: getEnv("POSTGRES_DATABASE", ""),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
