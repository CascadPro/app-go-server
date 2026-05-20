package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"cascade/config"
	"cascade/internal/models"
)

func GenerateRegisterToken(c *gin.Context) {
	user := models.User{
		ID:   uuid.New(),
		Role: models.RoleRegular,
	}

	config.DB.Model(&models.User{}).Create(&user)

	token := models.Token{
		ID:        uuid.New(),
		Token:     uuid.NewString(),
		Type:      models.TokenTypeRegister,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 12),
	}

	config.DB.Model(&models.Token{}).Create(&token)

	c.JSON(http.StatusOK, gin.H{"message": "Token successfully generated!"})
}
