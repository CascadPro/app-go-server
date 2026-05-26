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

	JwtSecretKey    string
	UploadSecretKey string

	S3Region          string
	S3Endpoint        string
	S3BucketName      string
	S3AccessKeyId     string
	S3SecretAccessKey string
	S3AllowedTags     string

	UseS3 bool
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logger.Error("No .env file found, using environment variables", nil)
	}

	applicationPort, err := strconv.Atoi(getEnv("APPLICATION_PORT", "4000"))
	if err != nil {
		logger.Error("❌ Failed to load config", err)
		return nil, err
	}

	postgresPort, err := strconv.Atoi(getEnv("POSTGRES_PORT", "5432"))
	if err != nil {
		logger.Error("❌ Failed to load config", err)
		return nil, err
	}

	redisPort, err := strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	if err != nil {
		logger.Error("❌ Failed to load config", err)
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
		PostgresDatabase: getEnv("POSTGRES_DB", ""),

		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisHost:     getEnv("REDIS_HOST", ""),
		RedisPort:     redisPort,

		JwtSecretKey:    getEnv("JWT_SECRET_KEY", ""),
		UploadSecretKey: getEnv("UPLOAD_SECRET_KEY", ""),

		S3Region:          getEnv("S3_REGION", ""),
		S3Endpoint:        getEnv("S3_ENDPOINT", ""),
		S3BucketName:      getEnv("S3_BUCKET_NAME", ""),
		S3AccessKeyId:     getEnv("S3_ACCESS_KEY_ID", ""),
		S3SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", ""),
		S3AllowedTags:     "docs,images,avatars",
		UseS3:             getEnv("S3_REGION", "") != "" && getEnv("S3_ENDPOINT", "") != "",
	}

	if config.UseS3 {
		if config.S3AccessKeyId == "" || config.S3SecretAccessKey == "" {
			logger.Error("❌ Missing S3 credentials: S3_ACCESS_KEY_ID and S3_SECRET_ACCESS_KEY must be set.", nil)

			return nil, err
		}
	} else {
		logger.Error("❌ S3 is not configured. Make sure S3_REGION and S3_ENDPOINT are set if you intend to use S3.", nil)
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
