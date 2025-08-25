package model

import "go-chat/internal/model"

type MessageDto struct {
	SendId  int64
	Message model.Message
}
