package tests

import (
	"go-chat/internal/model"
	"testing"
)

func TestId(t *testing.T) {

	messages := []*model.Message{
		{
			SenderId: 1,
		},
		{
			SenderId: 2,
		},
		{
			SenderId: 3,
		},
	}
	senderIds := make([]uint, len(messages))
	var id uint
	for i, msg := range messages {
		senderIds[i] = uint(msg.SenderId)
		id = uint(msg.SenderId)
	}
	t.Log(senderIds, id)
}
