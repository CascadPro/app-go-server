package auth

import (
	"net/http"

	"cascade/config"
	"cascade/internal/models"
	"cascade/pkg/filter"
	"cascade/pkg/logger"
	"cascade/pkg/utils"
	"cascade/pkg/utils/auth"

	"github.com/gin-gonic/gin"
)

type LoginDto struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64,alphanum"`
}

func Login(c *gin.Context) {
	var dto LoginDto
	if err := c.ShouldBind(&dto); err != nil {
		cause := err.Error()
		filter.Error(c, filter.ErrorParams{Status: http.StatusBadRequest, Message: "", Cause: &cause})
		return
	}

	var user models.User

	user_err := config.DB.Select("id", "role", "password").
		Where(&models.User{Email: &dto.Email}).
		First(&user).Error

	if user_err != nil {
		filter.Error(c, filter.ErrorParams{
			Status:  http.StatusBadRequest,
			Message: "Пользователя с такой эл. почтой не существует!"})
		return
	}

	if !user.CheckPassword(dto.Password) {
		filter.Error(c, filter.ErrorParams{Status: http.StatusUnauthorized, Message: "Неверный пароль!"})
		return
	}

	cfg, err := utils.LoadConfig()
	if err != nil {
		logger.Error("❌ Failed to load config", err)
		return
	}

	sessionID, session_err := auth.GenerateSessionID()
	if session_err != nil {
		logger.Error("Error during generating session ID!", session_err)
		return
	}

	if err := auth.CreateSession(sessionID, user.ID, auth.RefreshTokenLifetime); err != nil {
		logger.Error("Error during creating session!", err)
		return
	}

	accessToken, refreshToken := auth.IssueTokens(user.ID, user.Role, sessionID, cfg)

	c.SetCookie("refresh_token", refreshToken, int(auth.RefreshTokenLifetime), "/", cfg.Domain, false, true)

	filter.Success(c, "Вы успешно вошли в аккаунт!", gin.H{"access_token": accessToken})
}
