package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"cascade/pkg/logger"
	"cascade/pkg/utils"
)

type redis_db struct {
	DB  *redis.Client
	Ctx context.Context
}

var R redis_db

func ConnectRedisDatabase(cfg *utils.Config) {
	R.DB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
		Protocol: 2,
	})

	R.Ctx = context.Background()

	R.DB.ClusterAddSlots(R.Ctx, 1)

	logger.Info("✅ Redis database connection established successfully")
}
