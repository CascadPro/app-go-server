package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"cascade/config"
	"cascade/internal/handlers/auth"
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

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	{
		r_auth := router.Group("/auth")
		r_auth.GET("/token", auth.GenerateRegisterToken)
		r_auth.POST("/register", auth.Register)
		r_auth.POST("/login", auth.Login)
	}

	address := ":" + strconv.Itoa(cfg.ApplicationPort)

	logger.Info("🚀 Server is running at " + cfg.ApplicationUrl)

	if err := router.Run(address); err != nil {
		logger.Error("Error starting server", err)
	}
}
