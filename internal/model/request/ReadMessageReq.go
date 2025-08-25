package model

type ReadMessageReq struct {
	MessageId uint `json:"message_id"`
	UserId    uint `json:"user_id"`
}
