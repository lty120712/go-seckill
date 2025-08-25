package interfacesservice

import (
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	response "go-chat/internal/model/response"
)

// FriendServiceInterface 接口
type FriendServiceInterface interface {
	Add(id uint, friendIdList []uint) error
	ListReq(id uint) ([]response.FriendRequestVo, error)
	HandleReq(id int64, status model.Status) error
	Remove(userId uint, friendIdList []int64) error
	GroupCreate(id uint, req request.FriendGroupCreateRequest) error
	GroupDelete(groupId int64) error
	GroupUpdate(req request.FriendGroupUpdateRequest) error
	GroupList(userId uint) ([]response.GroupVo, error)
}
