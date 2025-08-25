package model

import (
	"encoding/json"
	"go-chat/internal/model"
)

// 发送消息请求结构体
type SendMessageRequest struct {
	SenderId   uint                  `json:"sender_id" gorm:"not null;comment:发送者ID"`    // 发送者ID（必填）
	ReceiverId *uint                 `json:"receiver_id" gorm:"comment:接收者ID（私聊使用）"`     // 接收者ID（仅用于私聊）
	GroupId    *uint                 `json:"group_id" gorm:"comment:群组ID（群聊使用）"`         // 群组ID（仅用于群聊）
	ReplyId    *uint                 `json:"reply_id" gorm:"comment:回复的消息ID"`            // 回复消息ID
	TargetType model.TargetType      `json:"target_type" gorm:"not null;comment:消息目标类型"` // 消息目标类型（0=私聊，1=群聊）
	Content    model.MessagePartList `json:"content" gorm:"type:json;comment:富文本消息内容"`   // 消息内容片段数组（JSON）
	Type       model.MessageType     `json:"type" gorm:"not null;comment:消息类型"`          // 消息类型（文本、图片、红包等）
	ExtraData  *json.RawMessage      `json:"extra_data" gorm:"type:json;comment:扩展字段"`   // 扩展字段（如红包、投票等结构）`
}
