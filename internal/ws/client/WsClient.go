package wsClient

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"go-chat/internal/utils/logUtil"
	"net/http"
	"sync"
)

type WebSocketManager struct {
	Server      *http.Server
	Upgrader    websocket.Upgrader
	Connections sync.Map // 存储所有连接，键为用户ID，值为WebSocket连接
}

var WebSocketClient *WebSocketManager

func (ws *WebSocketManager) Close() {
	if ws.Server != nil {
		if err := ws.Server.Close(); err != nil {
			logrus.Errorf("WebSocket 服务关闭失败: %s", err)
		}
	}
}
func (ws *WebSocketManager) SendMessageToOne(id int64, message interface{}) {
	conn, ok := WebSocketClient.Connections.Load(id)
	if !ok {
		logUtil.Warnf("没有找到用户 %s 的连接", id)
		return
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		logUtil.Errorf("消息序列化失败: %s", err)
		return
	}
	err = conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		logUtil.Errorf("向用户 %v 发送消息失败: %s", id, err)
	}
}
func (ws *WebSocketManager) SendMessageToMultiple(ids []int64, message interface{}) {
	for _, id := range ids {
		if conn, ok := WebSocketClient.Connections.Load(id); ok && conn != nil {
			ws.SendMessageToOne(id, message)
		}
	}
}

func (ws *WebSocketManager) SendMessageToAll(message interface{}) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		logUtil.Errorf("消息序列化失败: %s", err)
		return
	}
	WebSocketClient.Connections.Range(func(key, value interface{}) bool {
		conn := value.(*websocket.Conn)
		err := conn.WriteMessage(websocket.TextMessage, messageBytes)
		if err != nil {
			logUtil.Errorf("向用户 %v 发送消息失败: %s", key, err)
		}
		return true
	})
}

func (ws *WebSocketManager) GetOnlineUserIds() []int64 {
	var userIds []int64
	WebSocketClient.Connections.Range(func(key, value interface{}) bool {
		userIds = append(userIds, key.(int64))
		return true
	})
	return userIds
}
