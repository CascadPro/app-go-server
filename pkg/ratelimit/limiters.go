package ratelimit

import (
	"time"

	"cascade/config"
)

var AuthRateLimiter = NewRateLimiter(RateLimiterConfig{
	RedisClient: config.R.DB,
	Limit:       5,                // 5 запросов
	Window:      time.Minute * 15, // в 15 минут
	KeyPrefix:   "auth",
})

var AvatarRateLimiter = NewRateLimiter(RateLimiterConfig{
	RedisClient: config.R.DB,
	Limit:       10,               // 10 запросов
	Window:      time.Minute * 10, // в 10 минут
	KeyPrefix:   "avatar",
})
