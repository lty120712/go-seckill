package model

import "gorm.io/gorm"

type Friend struct {
	gorm.Model
	UserId   uint `json:"user_id"`   // 用户id
	FriendId uint `json:"friend_id"` // 好友id
}
