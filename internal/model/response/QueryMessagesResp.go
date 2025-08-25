package model

import (
	"go-chat/internal/model"
	"time"
)

// QueryMessagesResponse 查询消息列表
type QueryMessagesResponse struct {
	List    []*MessageVo `json:"list"`     // 消息列表
	Cursor  int64        `json:"cursor"`   // 下一页游标（最小 id）
	HasMore bool         `json:"has_more"` // 是否还有更多数据
}

type MessageVo struct {
	ID           uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	SenderId     int64
	ReceiverId   *int64
	GroupId      *int64
	ReplyId      *int64
	ReaderIdList *model.ReaderIdList
	TargetType   *model.TargetType
	Content      *model.MessagePartList
	Type         *model.MessageType
	Status       *model.Status
	ExtraData    interface{} `json:"extra_data" gorm:"type:json;comment:扩展字段"` // 扩展字段（如红包、投票等结构）

	//额外信息
	Reply              *MessageVo `json:"reply"`
	SenderNickName     *string
	SenderAvatar       *string
	SenderOnlineStatus *model.OnlineStatus
	IsRead             bool
}

func (m *MessageVo) GetFieldsFromMessage(msg *model.Message) {
	m.ID = msg.ID
	m.CreatedAt = msg.CreatedAt
	m.UpdatedAt = msg.UpdatedAt
	m.SenderId = msg.SenderId
	m.ReceiverId = msg.ReceiverId
	m.GroupId = msg.GroupId
	m.ReplyId = msg.ReplyId
	m.ReaderIdList = msg.ReaderIdList
	m.TargetType = msg.TargetType
	m.Content = msg.Content
	m.Type = msg.Type
	m.Status = msg.Status
	m.ExtraData = msg.ExtraData
}
