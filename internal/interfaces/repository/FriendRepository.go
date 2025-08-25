package interfaces

import (
	"go-chat/internal/model"
	response "go-chat/internal/model/response"
	"gorm.io/gorm"
)

type FriendRepositoryInterface interface {
	GetFriendsByUserId(id uint, tx ...*gorm.DB) ([]model.Friend, error)

	CreateFriendRequests(requests []model.FriendRequest, tx ...*gorm.DB) error

	BatchCreate(friends []model.Friend, tx ...*gorm.DB) error

	BatchDelete(userId uint, friendIds []uint, tx ...*gorm.DB) error
	BatchDeleteManyInverse(userId uint, friendIds []uint, tx ...*gorm.DB) error
	GetFriendsWithUserInfo(userId uint, tx ...*gorm.DB) ([]response.FriendVo, error)
}
