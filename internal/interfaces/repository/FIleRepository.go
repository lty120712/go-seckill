package interfaces

import (
	"go-chat/internal/model"
	"gorm.io/gorm"
)

type FileRepositoryInterface interface {
	Create(file *model.File, tx ...*gorm.DB) error
}
