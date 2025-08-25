package model

import (
	"go-chat/internal/model"
	"time"
)

// QueryMessagesRequest 查询消息列表请求参数
type QueryMessagesRequest struct {
	TargetId   uint              `json:"target_id" binding:"required"`   // 目标ID 好友id或群组id
	TargetType *model.TargetType `json:"target_type" binding:"required"` // 目标类型 私聊或群聊
	Cursor     uint              `json:"cursor"`                         // 游标 上次查询的最小id 默认为0
	Limit      int               `json:"limit"`                          // 限制数量
	//以下为拓展,暂时不会做
	MessageTypes *[]model.MessageType `json:"message_types"` // 只查特定类型的消息（如图片、文本）
	Keyword      *string              `json:"keyword"`       // 模糊搜索
	StartTime    time.Time            `json:"start_time"`    // 起始时间
	EndTime      time.Time            `json:"end_time"`      // 结束时间
}
