package media

import (
	"bytes"
	"cascade/config"
	"cascade/internal/repositories"
	"cascade/pkg/filter"
	"cascade/pkg/logger"
	"cascade/pkg/utils"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Fetch(c *gin.Context) {
	tag := c.Param("tag")
	id := c.Param("id")

	cfg, err := utils.LoadConfig()
	if err != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError})
		return
	}

	errMsg, err := utils.ValidateTagParam(tag, cfg)
	if err != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusBadRequest, Message: errMsg, Cause: err.Error()})
		return
	}

	validateErr := uuid.Validate(id)
	if validateErr != nil {
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

	s3Svc, s3Ctx, err := config.InitS3Session()
	if err != nil {
		logger.Error("❌ Failed to initialize AWS session", err)
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError})
		return
	}

	folder := utils.GetBucketFolder(tag)

	objectKey := fmt.Sprintf("%s%s", folder, fileRecord.ID)

	result, err := s3Svc.GetObject(s3Ctx, &s3.GetObjectInput{
		Bucket: aws.String(cfg.S3BucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		logger.Error("❌ S3 fetch error", err)
		filter.Error(c, filter.ErrorParams{
			Status:  http.StatusNotFound,
			Message: "Файл не найден!",
			Cause:   "Файла не существует в хранилище",
		})
		return
	}
	defer result.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(result.Body)
	if err != nil {
		logger.Error("❌ Failed to read file", err)
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError})
		return
	}

	fileBytes := buf.Bytes()
	contentType := *result.ContentType

	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", len(fileBytes)))
	c.Data(http.StatusOK, contentType, fileBytes)
}
