package auth

import (
	"net/http"
	"time"

	"cascade/config"
	"cascade/internal/models"
	"cascade/pkg/filter"
	"cascade/pkg/logger"

	"github.com/gin-gonic/gin"
)

type RegisterDto struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64,alphanum"`
	Token    string `json:"token" binding:"required,uuid"`
}

func Register(c *gin.Context) {
	var dto RegisterDto
	if err := c.ShouldBind(&dto); err != nil {
		filter.Error(c, filter.ErrorParams{
			Status:  http.StatusBadRequest,
			Message: "Неверно введены значения!",
			Cause:   err.Error(),
		})
		return
	}

	var token models.Token

	err := config.DB.Select("id", "user_id", "expires_at").
		Where(&models.Token{Token: dto.Token, Type: models.TokenTypeRegister}).
		First(&token).Error
	if err != nil {
		filter.Error(c, filter.ErrorParams{
			Status:  http.StatusBadRequest,
			Message: "Ключ неверный! Пожалуйста, запросите другой ключ",
		})
		return
	}

	if token.ExpiresAt.Before(time.Now()) {
		filter.Error(c, filter.ErrorParams{
			Status:  http.StatusBadRequest,
			Message: "Срок действия ключа истек! Пожалуйста, запросите новый ключ",
		})
		return
	}

	user := models.User{PasswordHash: &dto.Password}
	user.HashPassword()

	userErr := config.DB.
		Model(&models.User{}).
		Where(&models.User{ID: token.UserID}).
		Update("email", dto.Email).
		Update("password", user.PasswordHash).Error
	if userErr != nil {
		logger.Error("Can't update User model!", userErr)
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError})
		return
	}

	deleteErr := config.DB.Delete(&token).Error
	if deleteErr != nil {
		logger.Error("Can't delete register token!", deleteErr)
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError})
		return
	}

	filter.Success(c, "Вы успешно зарегистрировались!")
}
