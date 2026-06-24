package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"cascade/config"
	"cascade/internal/models"
	"cascade/pkg/filter"
)

type GenerateRegisterTokenDto struct {
	Name     string  `json:"name" binding:"required,max=255"`
	Surname  string  `json:"surname" binding:"required,max=255"`
	LastName *string `json:"last_name" binding:"omitempty,max=255"`
}

func GenerateRegisterToken(c *gin.Context) {
	var dto GenerateRegisterTokenDto
	if err := c.ShouldBind(&dto); err != nil {
		filter.Error(c, filter.ErrorParams{
			Status:  http.StatusBadRequest,
			Message: "Неверно введены значения!",
			Cause:   err.Error(),
		})
		return
	}

	user := models.User{
		ID:       uuid.New(),
		Role:     models.RoleRegular,
		Name:     dto.Name,
		Surname:  dto.Surname,
		LastName: dto.LastName,
	}

	config.DB.Model(&models.User{}).Create(&user)

	token := models.Token{
		ID:        uuid.New(),
		Token:     uuid.NewString(),
		Type:      models.TokenTypeRegister,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 12),
	}

	config.DB.Model(&models.Token{}).Create(&token)

	filter.Success(c, "Ключ успешно сгенерирован!", gin.H{"token": token.Token})
}
