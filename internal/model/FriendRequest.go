package model

import "gorm.io/gorm"

type FriendRequest struct {
	gorm.Model
	UserId   uint   `json:"user_id"`
	FriendId uint   `json:"friend_id"`
	Status   Status `json:"status"`
}
