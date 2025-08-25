package model

import "gorm.io/gorm"

// GroupAnnouncement 群组公告响应结构体
type GroupAnnouncement struct {
	gorm.Model
	GroupID   uint   `json:"group_id"`  // 群组ID
	Content   string `json:"content"`   // 公告内容
	Publisher uint   `json:"publisher"` // 发布者ID
}
