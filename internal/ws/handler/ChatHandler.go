package wsHandler

import (
	"go-chat/internal/db"
	"go-chat/internal/model"
	dto "go-chat/internal/model/dto"
	"go-chat/internal/utils/jsonUtil"
	"go-chat/internal/utils/redisUtil"
	wsClient "go-chat/internal/ws/client"
	"net/http"
	"strconv"
)

func (ws *WebSocketHandler) ChatHandler(sendId int64, data interface{}) {
	//个人限流
	key := "chat-limit:" + strconv.FormatInt(sendId, 10)
	isRateLimit := redisUtil.IsRateLimited(db.Redis, key, 99, 3)
	if isRateLimit {
		wsClient.WebSocketClient.SendMessageToOne(sendId, &model.Response{
			Code:    http.StatusTooManyRequests,
			Message: "limit message",
			Data:    nil,
		})
		return
	}
	bytes, err := jsonUtil.MarshalValue(data)
	if err != nil {
		wsClient.WebSocketClient.SendMessageToOne(sendId, &model.Response{
			Code:    http.StatusBadRequest,
			Message: "数据格式错误",
			Data:    nil,
		})
		return
	}

	var message = &model.Message{}
	if err := jsonUtil.UnmarshalValue(bytes, message); err != nil {
		wsClient.WebSocketClient.SendMessageToOne(sendId, &model.Response{
			Code:    http.StatusBadRequest,
			Message: "数据格式错误，无法反序列化成消息",
			Data:    nil,
		})
		return
	}
	message.InitFields()
	messageDto := dto.MessageDto{
		SendId:  sendId,
		Message: *message,
	}
	err = ws.rabbitMQManager.SendMessage("chat-exchange", "chat-storage", messageDto)
	if err != nil {
		wsClient.WebSocketClient.SendMessageToOne(sendId, &model.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
}
