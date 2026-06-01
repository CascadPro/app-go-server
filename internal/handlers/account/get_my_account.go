package account

import (
	"cascade/config"
	"cascade/internal/models"
	"cascade/pkg/filter"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMyAccount(c *gin.Context) {
	userID, isExists := c.Get("userID")
	if !isExists {
		filter.Error(c, filter.ErrorParams{
			Status: http.StatusInternalServerError,
			Cause:  "ID пользователя не был найден в контексте"})
		return
	}

	var user models.User

	err := config.DB.Select("id", "role", "email", "created_at", "updated_at").First(&user, "id = ?", userID).Error
	if err != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusNotFound, Message: "Пользователь не был найден!"})
		return
	}

	c.JSON(http.StatusOK, user)
}
