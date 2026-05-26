package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"cascade/config"
	"cascade/pkg/filter"
	"cascade/pkg/utils/authutils/sessions"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiterConfig struct {
	RedisClient *redis.Client
	Limit       int64         // максимальное количество запросов
	Window      time.Duration // временной интервал
	KeyPrefix   string        // префикс для ключей в Redis
}

type RateLimiter struct {
	config RateLimiterConfig
}

func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	return &RateLimiter{
		config: config,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Можно также использовать user ID или другой идентификатор
		// clientIP := c.GetHeader("X-Forwarded-For")

		key := fmt.Sprintf("%s%s:%s", sessions.RedisRateLimitFolder, rl.config.KeyPrefix, clientIP)

		allowed, remaining, resetTime, err := rl.isAllowed(key)
		if err != nil {
			filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError})
			return
		}

		c.Header("X-RateLimit-Limit", strconv.FormatInt(rl.config.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

		if !allowed {
			filter.Error(c, filter.ErrorParams{
				Status:  http.StatusTooManyRequests,
				Message: "Too Many Requests",
				Cause:   "Rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) isAllowed(key string) (bool, int64, int64, error) {
	luaScript := `
		local key = KEYS[1]
		local limit = tonumber(ARGV[1])
		local window = tonumber(ARGV[2])
		
		local current = redis.call('GET', key)
		
		if current == false then
			redis.call('SET', key, 1)
			redis.call('EXPIRE', key, window)
			return {1, limit - 1, window}
		end
		
		local count = tonumber(current)
		
		if count >= limit then
			local ttl = redis.call('TTL', key)
			return {0, 0, ttl}
		end
		
		local new_count = redis.call('INCR', key)
		local ttl = redis.call('TTL', key)
		
		if ttl == -1 then
			redis.call('EXPIRE', key, window)
			ttl = window
		end
		
		return {1, limit - new_count, ttl}
	`

	result, err := rl.config.RedisClient.Eval(config.R.Ctx, luaScript, []string{key},
		rl.config.Limit, int64(rl.config.Window.Seconds())).Result()

	if err != nil {
		return false, 0, 0, err
	}

	resultSlice := result.([]any)
	allowed := resultSlice[0].(int64) == 1
	remaining := resultSlice[1].(int64)
	ttl := resultSlice[2].(int64)

	resetTime := time.Now().Add(time.Duration(ttl) * time.Second).Unix()

	return allowed, remaining, resetTime, nil
}
