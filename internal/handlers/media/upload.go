package media

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"

	"cascade/config"
	"cascade/internal/models"
	"cascade/internal/repositories"
	"cascade/pkg/filter"
	"cascade/pkg/logger"
	"cascade/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Upload(c *gin.Context) {
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

	file, err := getFileFromRequest(c)
	if err != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError, Cause: err.Error()})
		return
	}

	buf, err := readFile(file)
	if err != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError, Cause: err.Error()})
		return
	}

	tag := c.PostForm("tag")

	errMsg, err := utils.ValidateTagParam(tag, cfg)
	if err != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusBadRequest, Message: errMsg, Cause: err.Error()})
		return
	}

	fileID, err := uuid.NewUUID()
	if err != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError, Cause: err.Error()})
		return
	}

	s3Svc, s3Ctx, err := config.InitS3Session()
	if err != nil {
		logger.Error("❌ Failed to initialize AWS session", err)
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError})
		return
	}

	folder := utils.GetBucketFolder(tag)

	_, s3Err := s3Svc.PutObject(s3Ctx, &s3.PutObjectInput{
		Bucket:      aws.String(cfg.S3BucketName),
		Key:         aws.String(folder + fileID.String()),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})
	if s3Err != nil {
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError})
		return
	}

	fileRepo := repositories.NewFileRepository()

	fileRecord := models.File{
		ID:          fileID,
		Tag:         tag,
		Filename:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Size:        int(file.Size),
		Deleted:     nil,
	}

	if err := fileRepo.CreateFile(&fileRecord); err != nil {
		logger.Error("❌ Failed to save file metadata in database", err)
		filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError})
		return
	}

	filter.Success(c, "Файл успешно загружен!", gin.H{
		"file_id":  fileID,
		"filename": file.Filename,
	})
}

func getFileFromRequest(c *gin.Context) (*multipart.FileHeader, error) {
	file, err := c.FormFile("file")
	if err != nil {
		logger.Error("❌ Failed to get file from request", err)
		return nil, fmt.Errorf("не получилось извлечь файл из запроса: %w", err)
	}
	return file, nil
}

func readFile(file *multipart.FileHeader) (*bytes.Buffer, error) {
	src, err := file.Open()
	if err != nil {
		logger.Error("❌ Failed to open a file!", err)
		return nil, fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer src.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(src)
	if err != nil {
		logger.Error("❌ Failed to read a file!", err)
		return nil, fmt.Errorf("не удалось прочитать файл: %w", err)
	}
	return buf, nil
}
