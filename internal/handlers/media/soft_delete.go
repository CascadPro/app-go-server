package media

import (
	"cascade/internal/repositories"
	"cascade/pkg/filter"
	"cascade/pkg/logger"
	"cascade/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SoftDelete(c *gin.Context) {
	cfg, err := utils.LoadConfig()
	if err != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError})
		return
	}

	secretKey := c.GetHeader("X-Upload-Secret")
	if secretKey != cfg.UploadSecretKey {
		filter.Error(c, filter.ErrorParams{
			Status:  http.StatusForbidden,
			Message: "Неправильный секретный ключ!",
		})
		return
	}

	tag := c.Param("tag")
	id := c.Param("id")

	errMsg, err := utils.ValidateTagParam(tag, cfg)
	if err != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusBadRequest, Message: errMsg, Cause: err.Error()})
		return
	}

	id_err := uuid.Validate(id)
	if id_err != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusBadRequest, Message: "Неверный формат ID! (UUID)"})
		return
	}

	fileRepo := repositories.NewFileRepository()

	fileRecord, err := fileRepo.GetFileByID(id, tag)
	if err != nil || fileRecord == nil {
		filter.Error(c, filter.ErrorParams{
			Status:  http.StatusNotFound,
			Message: "Файл не найден!",
			Cause:   "Записи о файле не существует!",
		})
		return
	}

	err = fileRepo.SoftDeleteFile(id, tag)
	if err != nil {
		logger.Error("❌ Failed to mark a file as deleted!", err)
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError, Message: "Не удалось удалить файл!"})
		return
	}

	filter.Success(c, "Файл успешно помечен как удаленный")
}
