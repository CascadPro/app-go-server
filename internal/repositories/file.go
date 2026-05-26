package repositories

import (
	"cascade/config"
	"cascade/internal/models"

	"gorm.io/gorm"
)

type FileRepository struct{}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (r *FileRepository) GetFileByID(id, tag string) (*models.File, error) {
	var file models.File
	if err := config.DB.First(&file, "id = ? AND tag = ? AND deleted = ?", id, tag, false).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}
