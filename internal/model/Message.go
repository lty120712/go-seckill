package model

import (
	"database/sql/driver"
	"go-chat/internal/utils/jsonUtil"
	"gorm.io/gorm"
)

// 消息结构体
type Message struct {
	gorm.Model
	SenderId     int64            `json:"sender_id" gorm:"not null;comment:发送者ID"`          // 发送者ID（必填）
	ReceiverId   *int64           `json:"receiver_id" gorm:"comment:接收者ID（私聊使用）"`           // 接收者ID（仅用于私聊）
	GroupId      *int64           `json:"group_id" gorm:"comment:群组ID（群聊使用）"`               // 群组ID（仅用于群聊）
	ReplyId      *int64           `json:"reply_id" gorm:"comment:回复的消息ID"`                  // 回复消息ID
	ReaderIdList *ReaderIdList    `json:"reader_id_list" gorm:"type:json;comment:已读用户ID列表"` // 已读用户ID数组，JSON 存储
	TargetType   *TargetType      `json:"target_type" gorm:"not null;comment:消息目标类型"`       // 消息目标类型（0=私聊，1=群聊）
	Content      *MessagePartList `json:"content" gorm:"type:json;comment:富文本消息内容"`         // 消息内容片段数组（JSON）
	Type         *MessageType     `json:"type" gorm:"not null;comment:消息类型"`                // 消息类型（文本、图片、红包等）
	Status       *Status          `json:"status" gorm:"not null;comment:消息状态"`              // 消息状态（0=撤回，1=正常）
	ExtraData    interface{}      `json:"extra_data" gorm:"type:json;comment:扩展字段"`         // 扩展字段（如红包、投票等结构）
}

func (m *Message) TableName() string {
	return "messages"
}

// InitFields 修复数据
func (m *Message) InitFields() {
	// 处理 Status 字段
	if m.Status == nil {
		m.Status = new(Status) // 为 Status 分配内存
		*m.Status = Enable     // 然后赋值
	}
	// 处理 ReaderIdList 字段
	if m.ReaderIdList == nil {
		m.ReaderIdList = new(ReaderIdList)
		*m.ReaderIdList = make(ReaderIdList, 0)
	}
}

type MessageType int

const (
	TextContent MessageType = iota // 聊天消息

	ImageContent // 图片

	VoiceContent // 语音

	RedBagContent //红包

	ForwardedCotent //转发
	SystemContent   //系统消息
)

type TargetType int

const (
	PrivateTarget TargetType = iota // 私聊
	GroupTarget                     // 群聊
)

type ContentType string

const (
	Text  ContentType = "text"  // 文本
	Emoji ContentType = "emoji" // 表情
	Image ContentType = "image" // 图片
	Link  ContentType = "link"  // 链接
)

type MessagePart struct {
	Type    ContentType `json:"type"`    // 内容类型（text, emoji, image, link）
	Content *string     `json:"content"` // 内容（如文本、图片 URL、链接等）
}

type MessagePartList []*MessagePart

func (parts *MessagePartList) Value() (driver.Value, error) {
	return jsonUtil.MarshalValue(parts)
}

func (parts *MessagePartList) Scan(value interface{}) error {
	return jsonUtil.UnmarshalValue(value, parts)
}

type ReaderIdList []uint

func (ids *ReaderIdList) Value() (driver.Value, error) {
	return jsonUtil.MarshalValue(ids)
}
func (ids *ReaderIdList) Scan(value interface{}) error {
	return jsonUtil.UnmarshalValue(value, ids)
}
