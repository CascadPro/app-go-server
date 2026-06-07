package account_avatar

import (
	"net/http"
	"path"

	"cascade/config"
	"cascade/internal/models"
	"cascade/pkg/filter"
	"cascade/pkg/logger"
	"cascade/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DeleteAvatar(c *gin.Context) {
	cfg, err := utils.LoadConfig()
	if err != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError})
		return
	}

	userID, isExists := c.Get("userID")
	if !isExists {
		filter.Error(c, filter.ErrorParams{
			Status: http.StatusInternalServerError,
			Cause:  "ID пользователя не был найден в контексте"})
		return
	}

	var user models.User

	userErr := config.DB.Select("id", "avatar").First(&user, "id = ?", userID).Error
	if userErr != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusNotFound, Message: "Пользователь не был найден!"})
		return
	}

	if err := uuid.Validate(*user.Avatar); err != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusConflict, Message: "Нельзя удалить стандартное фото!"})
		return
	}

	response, err := http.NewRequest(
		http.MethodDelete, path.Join(cfg.ApplicationUrl, "media", utils.TagAvatars, *user.Avatar), http.NoBody,
	)
	if err != nil || response.Response.StatusCode != 200 {
		filter.Error(c, filter.ErrorParams{
			Status:  response.Response.StatusCode,
			Message: "Произошла ошибка во время удаления аватара!",
			Cause:   err.Error(),
		})
		return
	}

	user.Avatar = nil

	updateErr := config.DB.Save(&user).Error
	if updateErr != nil {
		logger.Error("❌ Failed to update user's avatar info!", updateErr)
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError})
		return
	}

	filter.Success(c, "Аватар успешно удален!")
}
