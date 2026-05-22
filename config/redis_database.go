package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"cascade/pkg/logger"
	"cascade/pkg/utils"

)

var Redis *redis.Client
var RedisCtx context.Context

func ConnectRedisDatabase(cfg *utils.Config) {
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
		Protocol: 2,
	})

	RedisCtx = context.Background()

	Redis.ClusterAddSlots(RedisCtx, 1)

	logger.Info("✅ Redis database connection established successfully")
}
