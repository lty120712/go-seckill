// Package interfaces
package interfaces

import (
	"go-chat/internal/model"
	response "go-chat/internal/model/response"
	"gorm.io/gorm"
)

// UserRepositoryInterface 用户仓库接口
type UserRepositoryInterface interface {
	GetById(id uint, tx ...*gorm.DB) (user *model.User, err error)
	GetByName(username *string, tx ...*gorm.DB) (user *model.User, err error)
	Save(user *model.User, tx ...*gorm.DB) (err error)
	UpdateFields(id uint, updates map[string]interface{}, tx ...*gorm.DB) error
	GetNickNamesByIds(ids []uint, tx ...*gorm.DB) (map[uint]string, error)
	GetNickNamesById(id uint, tx ...*gorm.DB) (nickname string, err error)
	GetByIdList(userIdList []uint, tx ...*gorm.DB) (userList []model.User, err error)
	GetVoById(id uint, tx ...*gorm.DB) (userVo response.UserVO, err error)
	UpdateHeartbeatTime(userId int64, time int64, tx ...*gorm.DB) error

	GetUsersWithHeartbeatBefore(cutoffTime int64, tx ...*gorm.DB) ([]model.User, error)
}
