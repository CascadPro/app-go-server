package utils

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"cascade/pkg/logger"
)

type Config struct {
	Domain          string
	ApplicationPort int
	ApplicationUrl  string
	AllowedOrigins  []string

	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string

	RedisPassword string
	RedisHost     string
	RedisPort     int

	JwtSecretKey string
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

	redisPort, err := strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	if err != nil {
		return nil, err
	}

	allowedOrigins := os.ExpandEnv(getEnv("ALLOWED_ORIGINS", ""))
	allowedOriginsList := strings.Split(allowedOrigins, ",")

	config := &Config{
		Domain:          getEnv("DOMAIN", "localhost"),
		ApplicationPort: applicationPort,
		ApplicationUrl:  os.ExpandEnv(getEnv("APPLICATION_URL", "")),
		AllowedOrigins:  allowedOriginsList,

		PostgresUser:     getEnv("POSTGRES_USER", ""),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
		PostgresHost:     getEnv("POSTGRES_HOST", ""),
		PostgresPort:     postgresPort,
		PostgresDatabase: getEnv("POSTGRES_DATABASE", ""),

		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisHost:     getEnv("REDIS_HOST", ""),
		RedisPort:     redisPort,

		JwtSecretKey: getEnv("JWT_SECRET_KEY", ""),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
