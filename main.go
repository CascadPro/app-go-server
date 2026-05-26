package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"cascade/config"
	"cascade/internal/handlers/account"
	"cascade/internal/handlers/auth"
	"cascade/internal/middlewares"
	"cascade/internal/models"
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

	_, _, s3_err := config.InitS3Session()
	if s3_err != nil {
		logger.Error("❌ Failed to initialize AWS session", s3_err)
		return
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Media Content
	{
		r := router.Group("/media")
		r.Use(middlewares.MediaCors(cfg))

		router.GET("/:tag/:id", media.Fetch)
	}

	// Auth group
	authRaLm := middlewares.NewRateLimiter(middlewares.RateLimiterConfig{
		RedisClient: config.R.DB,
		Limit:       5,                // 5 запросов
		Window:      time.Minute * 15, // в 15 минут
		KeyPrefix:   "auth",
	})

	{
		r := router.Group("/auth")
		r.Use(authRaLm.Middleware())

		r.POST("/login", auth.Login)
		r.GET("/login/refresh", auth.GetNewTokens)
		r.POST("/logout", middlewares.AuthMiddleware(), auth.Logout)

		r.POST("/register", auth.Register)
		r.GET("/register/token", middlewares.AuthMiddleware(models.RoleAdmin, models.RoleDirector),
			auth.GenerateRegisterToken)
	}

	// Account group
	accountRaLm := middlewares.NewRateLimiter(middlewares.RateLimiterConfig{
		RedisClient: config.R.DB,
		Limit:       10,          // 10 запросов
		Window:      time.Minute, // в минуту
		KeyPrefix:   "account",
	})

	{
		r := router.Group("/account")
		r.Use(accountRaLm.Middleware())
		r.Use(middlewares.AuthMiddleware())

		r.GET("/my", account.GetMyAccount)
	}

	address := ":" + strconv.Itoa(cfg.ApplicationPort)

	logger.Info("🚀 Server is running at " + cfg.ApplicationUrl)

	if err := router.Run(address); err != nil {
		logger.Error("Error starting server", err)
	}
}
