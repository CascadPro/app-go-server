package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"cascade/config"
	"cascade/internal/handlers/account"
	account_avatar "cascade/internal/handlers/account/avatar"
	"cascade/internal/handlers/auth"
	"cascade/internal/handlers/media"
	"cascade/internal/middlewares"
	"cascade/internal/models"
	"cascade/pkg/logger"
	"cascade/pkg/ratelimit"
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

	if cfg.Connection == "online" {
		_, _, s3_err := config.InitS3Session()
		if s3_err != nil {
			logger.Error("❌ Failed to initialize AWS session", s3_err)
			return
		}
		logger.Info("✅ AWS session is initialized successfully")
	} else {
		logger.Info("ℹ	Skipped initializing AWS session, cause connection is offline")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Media Content
	{
		r := router.Group("/media")
		r.Use(middlewares.MediaCors(cfg))

		r.GET("/:tag/:id", media.Fetch)
		r.POST("/upload", media.Upload)
		r.DELETE("/:tag/:id", media.SoftDelete)
	}

	// Auth group
	{
		r := router.Group("/auth")
		r.Use(ratelimit.AuthRateLimiter.Middleware())

		r.POST("/login", auth.Login)
		r.GET("/login/refresh", auth.GetNewTokens)
		r.POST("/logout", middlewares.AuthMiddleware(), auth.Logout)

		r.POST("/register", auth.Register)
		r.POST("/register/token", middlewares.AuthMiddleware(models.RoleAdmin, models.RoleDirector),
			auth.GenerateRegisterToken)
	}


	// Account group
	{
		r := router.Group("/account")
		r.Use(middlewares.AuthMiddleware())

		r.GET("/my", account.GetMyAccount)

		{
			ar := r.Group("/avatar")
			ar.Use(ratelimit.AvatarRateLimiter.Middleware())

			ar.PATCH("/update", account_avatar.UpdateAvatar)
			ar.DELETE("/delete", account_avatar.DeleteAvatar)
		}
	}

	address := ":" + strconv.Itoa(cfg.ApplicationPort)

	logger.Info("⚙	Connection mode is " + cfg.Connection)
	logger.Info("🚀 Server is running at " + cfg.ApplicationUrl)

	if err := router.Run(address); err != nil {
		logger.Error("Error starting server", err)
	}
}
