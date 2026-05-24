package auth

import (
	"net/http"
	"time"

	"cascade/config"
	"cascade/internal/models"
	"cascade/pkg/filter"
	"cascade/pkg/logger"
	"cascade/pkg/utils"
	"cascade/pkg/utils/authutils"
	"cascade/pkg/utils/authutils/sessions"

	"github.com/gin-gonic/gin"
)

type LoginDto struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64,alphanum"`
}

func Login(c *gin.Context) {
	var dto LoginDto
	if err := c.ShouldBind(&dto); err != nil {
		filter.Error(c, filter.ErrorParams{
			Status:  http.StatusBadRequest,
			Message: "Неверно введены значения!",
			Cause:   err.Error()})
	}

	var user models.User

	err := config.DB.Select("id", "role", "password").
		Where(&models.User{Email: &dto.Email}).
		First(&user).Error

	if err != nil {
		filter.Error(c, filter.ErrorParams{
			Status:  http.StatusBadRequest,
			Message: "Пользователя с такой эл. почтой не существует!"})
	}

	if !user.CheckPassword(dto.Password) {
		filter.Error(c, filter.ErrorParams{Status: http.StatusUnauthorized, Message: "Неверный пароль!"})
	}

	cfg, err := utils.LoadConfig()
	if err != nil {
		filter.Error(c, filter.ErrorParams{
			Status:  http.StatusInternalServerError,
			Message: "Something went wrong!",
			Cause:   err.Error()})
	}

	sessionID, err := sessions.GenerateSessionID()
	if err != nil {
		logger.Error("Error during generating session ID!", err)
		return
	}

	if err := sessions.CreateSession(c.Request, sessionID, user.ID, authutils.RefreshTokenLifetime); err != nil {
		logger.Error("Error during creating session!", err)
		return
	}

	accessToken, refreshToken := authutils.IssueTokens(user.ID, user.Role, sessionID, cfg)

	c.SetCookieData(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Domain:   cfg.Domain,
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(authutils.RefreshTokenLifetime),
	})

	filter.Success(c, "Вы успешно вошли в аккаунт!", gin.H{"access_token": accessToken})
}
