package interfaces

import (
	"go-chat/internal/model"
	"gorm.io/gorm"
)

type FriendGroupRepositoryInterface interface {
	CreateGroup(group *model.FriendGroup, tx ...*gorm.DB) error
	GetGroupsByUserId(userId uint, tx ...*gorm.DB) ([]model.FriendGroup, error)
	UpdateGroup(group *model.FriendGroup, tx ...*gorm.DB) error
	DeleteGroupById(groupId uint, tx ...*gorm.DB) error
	GetGroupById(groupId uint, tx ...*gorm.DB) (*model.FriendGroup, error)
}
