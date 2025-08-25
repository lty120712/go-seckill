package wsHandler

import (
	"go-chat/internal/model"
	wsClient "go-chat/internal/ws/client"
	wsMessage "go-chat/internal/ws/message"
	"net/http"
	"time"
)

func (ws *WebSocketHandler) HeartBeatHandler(sendId int64, _ interface{}) {
	timestamp := time.Now().Unix()
	err := ws.userService.UpdateHeartbeatTime(sendId, timestamp)
	if err != nil {
		wsClient.WebSocketClient.SendMessageToOne(sendId, &model.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	// 返回心跳确认
	wsClient.WebSocketClient.SendMessageToOne(sendId, &model.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data: &wsMessage.Message{
			SendId: sendId,
			Type:   wsMessage.HeartBeatAck,
			Data:   nil,
			Time:   time.Now(),
		},
	})
}
