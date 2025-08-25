package interfaces

import (
	"go-chat/internal/model"
	response "go-chat/internal/model/response"
	"gorm.io/gorm"
)

type FriendRequestRepositoryInterface interface {
	GetFriendRequestsByUser(userId uint, tx ...*gorm.DB) ([]model.FriendRequest, error)
	GetSentFriendRequests(userId uint, tx ...*gorm.DB) ([]response.FriendRequestVo, error)
	GetReceivedFriendRequests(userId uint, tx ...*gorm.DB) ([]response.FriendRequestVo, error)
	GetById(id int64, tx ...*gorm.DB) (*model.FriendRequest, error)

	UpdateStatus(id int64, status model.Status, tx ...*gorm.DB) error
}
