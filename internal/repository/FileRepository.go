package repository

import (
	"go-chat/internal/db"
	"go-chat/internal/model"
	"gorm.io/gorm"
)

type FileRepository struct {
}

var (
	FileRepositoryInstance *FileRepository
)

func InitFileRepository() {
	FileRepositoryInstance = &FileRepository{}
}

func (r *FileRepository) Create(file *model.File, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Create(file).Error
}
