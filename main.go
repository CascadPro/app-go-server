package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"cascade/config"
	"cascade/internal/handlers/auth"
	"cascade/internal/middlewares"
	"cascade/pkg/logger"
	"cascade/pkg/utils"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	logger.InitLogger()

	cfg, err := utils.LoadConfig()
	if err != nil {
		return
	}

	config.ConnectPgDatabase(cfg)
	config.ConnectRedisDatabase(cfg)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/ping", middlewares.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"ms":      c.GetTime("").Nanosecond(),
		})
	})

	// Auth group
	authRaLm := middlewares.NewRateLimiter(middlewares.RateLimiterConfig{
		RedisClient: config.Redis,
		Limit:       5,                // 10 запросов
		Window:      time.Minute * 15, // в минуту
		KeyPrefix:   "auth",           // префикс для ключей
	})

	{
		r_auth := router.Group("/auth")
		r_auth.Use(authRaLm.Middleware())

		r_auth.GET("/token", auth.GenerateRegisterToken)
		r_auth.GET("/login/refresh", auth.GetNewTokens)
		r_auth.POST("/register", auth.Register)
		r_auth.POST("/login", auth.Login)
	}

	address := ":" + strconv.Itoa(cfg.ApplicationPort)

	logger.Info("🚀 Server is running at " + cfg.ApplicationUrl)

	if err := router.Run(address); err != nil {
		logger.Error("Error starting server", err)
	}
}
