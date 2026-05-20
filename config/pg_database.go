package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"cascade/internal/models"
	"cascade/pkg/logger"
	"cascade/pkg/utils"
)

var DB *gorm.DB

func ConnectPgDatabase(cfg *utils.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
		cfg.PostgresPort,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("❌ Failed to establish PostgreSQL database connection: %v", err)
	}

	logger.Info("✅ PostgreSQL database connection established successfully")

	err = DB.AutoMigrate(&models.User{})

	if err != nil {
		logger.Error("❌ Error during PostgreSQL migration: %v", err)
	}

	logger.Info("✅ PostgreSQL database migrated successfully")
}
