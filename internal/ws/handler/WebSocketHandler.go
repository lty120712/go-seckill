package wsHandler

import (
	interfacemanager "go-chat/internal/interfaces/manager"
	interfacesservice "go-chat/internal/interfaces/service"
	"go-chat/internal/model"
	"go-chat/internal/utils/jsonUtil"
	wsClient "go-chat/internal/ws/client"
	wsMessage "go-chat/internal/ws/message"
	"net/http"
)

type WebSocketHandler struct {
	userService     interfacesservice.UserServiceInterface
	messageService  interfacesservice.MessageServiceInterface
	groupService    interfacesservice.GroupServiceInterface
	rabbitMQManager interfacemanager.RabbitMQManager
}

var (
	WebSocketHandlerInstance *WebSocketHandler
)

func InitWebSocketHandler(userService interfacesservice.UserServiceInterface,
	messageService interfacesservice.MessageServiceInterface,
	groupService interfacesservice.GroupServiceInterface,
	rabbitMQManager interfacemanager.RabbitMQManager) {
	WebSocketHandlerInstance = &WebSocketHandler{
		userService:     userService,
		messageService:  messageService,
		groupService:    groupService,
		rabbitMQManager: rabbitMQManager,
	}
}
func (ws *WebSocketHandler) MessageHandler(id int64, messageBytes []byte) {
	// 处理消息
	message := &wsMessage.Message{}
	err := jsonUtil.UnmarshalValue(messageBytes, message)
	if err != nil {
		wsClient.WebSocketClient.SendMessageToOne(id, &model.Response{
			Code:    http.StatusBadRequest,
			Message: "数据格式错误",
			Data:    nil,
		})
		return
	}
	switch message.Type {
	case wsMessage.Chat:
		ws.ChatHandler(message.SendId, message.Data)
	case wsMessage.HeartBeat:
		ws.HeartBeatHandler(message.SendId, message.Data)
	default:
		wsClient.WebSocketClient.SendMessageToOne(id, &model.Response{
			Code:    http.StatusBadRequest,
			Message: "数据格式错误",
			Data:    nil,
		})
	}
}
