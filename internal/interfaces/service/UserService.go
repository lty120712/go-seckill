package interfacesservice

import (
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	response "go-chat/internal/model/response"
)

// UserServiceInterface 接口
type UserServiceInterface interface {
	Register(username, password, rePassword *string) (err error)
	Login(username, password *string) (token string, err error)
	Logout(id uint)

	OnlineStatusChange(id uint, onlineStatus model.OnlineStatus) error
	UpdateUser(updateRequest *request.UserUpdateRequest) error
	GetUserInfo(id uint) (response.UserVO, error)

	UpdateHeartbeatTime(userId int64, time int64) error

	CheckOfflineUsers() error
}
