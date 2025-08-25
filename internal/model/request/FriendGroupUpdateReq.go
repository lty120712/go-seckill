package model

import "go-chat/internal/model"

type FriendGroupUpdateRequest struct {
	GroupId      uint               `json:"group_id" binding:"required"`
	Name         *string            `json:"name"`           // 可选，修改分组名称
	FriendIdList model.FriendIdList `json:"friend_id_list"` // 替换好友列表，允许空数组
}
