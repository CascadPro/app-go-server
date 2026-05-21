package auth

import (
	"net/http"
	"time"

	"cascade/config"
	"cascade/internal/models"
	"cascade/pkg/filter"
	"cascade/pkg/logger"
	"cascade/pkg/utils"

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
		cause := err.Error()
		filter.Error(c, filter.ErrorParams{Status: http.StatusBadRequest, Message: "", Cause: &cause})
		return
	}

	var token models.Token

	token_err := config.DB.Select("id", "user_id", "expires_at").
		Where(&models.Token{Token: dto.Token, Type: models.TokenTypeRegister}).
		First(&token).Error
	if token_err != nil {
		filter.Error(c, filter.ErrorParams{
			Status:  http.StatusBadRequest,
			Message: "Ключ неверный! Пожалуйста, запросите другой ключ"})
		return
	}

	if token.ExpiresAt.Before(time.Now()) {
		filter.Error(c, filter.ErrorParams{
			Status:  http.StatusBadRequest,
			Message: "Срок действия ключа истек! Пожалуйста, запросите новый ключ"})
		return
	}

	password_hash, hash_err := utils.GenerateHash(dto.Password)
	if hash_err != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError})
		logger.Error("Error while generating hash string!", hash_err)
		return
	}

	user_err := config.DB.
		Model(&models.User{}).
		Where(&models.User{ID: token.UserID}).
		Update("email", dto.Email).
		Update("password", password_hash).Error
	if user_err != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError})
		logger.Error("Can't update User model!", user_err)
		return
	}

	delete_err := config.DB.Delete(&token).Error
	if delete_err != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError})
		logger.Error("Can't delete register token!", delete_err)
		return
	}

	filter.Success(c, "Вы успешно зарегистрировались!")
}
