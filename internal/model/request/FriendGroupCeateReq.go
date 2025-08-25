package model

import "go-chat/internal/model"

type FriendGroupCreateRequest struct {
	Name         string              `json:"name" binding:"required"`
	FriendIdList *model.FriendIdList `json:"friend_id_list" `
}
