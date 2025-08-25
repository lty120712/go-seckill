package interfaces

import (
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
)

type MessageRepositoryInterface interface {
	Save(message *model.Message) (err error)
	GetById(id uint) (message *model.Message, err error)
	UpdateFields(id uint, fields map[string]interface{}) (err error)
	QueryHistoryMessages(userId uint, req *request.QueryMessagesRequest) ([]*model.Message, error)
}
