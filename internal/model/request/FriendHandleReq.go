package model

import "go-chat/internal/model"

type FriendHandlerReq struct {
	Id     int64        `json:"id"`     //记录id
	Status model.Status `json:"status"` //状态
}
