package interfacesservice

import (
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	response "go-chat/internal/model/response"
)

type MessageServiceInterface interface {
	// SendMessage 发送消息（支持私聊和群聊）
	// msg 是已经构造好的 message 对象（建议外部构建 content 等）
	SendMessage(msg *model.Message) (*response.MessageVo, error)

	GetMessageById(id uint) (*response.MessageVo, error)

	ReadMessage(messageId uint, userId uint) error

	QueryMessages(userId uint, req *request.QueryMessagesRequest) (*response.QueryMessagesResponse, error)
	Revoke(userId uint, messageId uint) error
}
