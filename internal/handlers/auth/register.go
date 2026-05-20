package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"cascade/config"
	"cascade/internal/models"
	"cascade/pkg/logger"
	"cascade/pkg/utils"
)

type RegisterDto struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
	Token    string `json:"token" binding:"required,uuid"`
}

func Register(c *gin.Context) {
	var dto RegisterDto
	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var token models.Token

	token_err := config.DB.Model(&models.Token{}).
		Select("id", "user_id", "expires_at").
		Where(&models.Token{Token: dto.Token, Type: models.TokenTypeRegister}).
		First(&token).Error
	if token_err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Token is invalid! Please, ask another token"})
		return
	}

	if token.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token is expired! Please, ask new token"})
		return
	}

	password_hash, hash_err := utils.GenerateHash(dto.Password)
	if hash_err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong!"})
		logger.Error("Error while generating hash string!", hash_err)
		return
	}

	user_err := config.DB.Model(&models.User{}).
		Where(&models.User{ID: token.UserID}).
		Updates(map[string]interface{}{"email": dto.Email, "password": password_hash}).Error
	if user_err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong!"})
		logger.Error("Can't update User model!", user_err)
		return
	}

	delete_err := config.DB.Model(&models.Token{}).Delete(map[string]interface{}{"id": token.ID}).Error
	if delete_err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong!"})
		logger.Error("Can't delete register token!", delete_err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you have registered!"})
}
