package manager

import (
	"github.com/gorilla/websocket"
	"go-chat/configs"
	"go-chat/internal/utils/logUtil"
	wsClient "go-chat/internal/ws/client"
	"go-chat/internal/ws/handler"
	"net/http"
	"strconv"
)

// InitWebSocket 初始化 WebSocket
func InitWebSocket() {
	webSocketConfig := configs.AppConfig.WebSocket
	wsClient.WebSocketClient = &wsClient.WebSocketManager{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	// 监听客户端连接
	http.HandleFunc("/ws", handleWebSocket)

	// 启动 WebSocket 服务
	go func() {
		wsClient.WebSocketClient.Server = &http.Server{
			Addr: webSocketConfig.Addr,
		}
		if err := wsClient.WebSocketClient.Server.ListenAndServe(); err != nil {
			logUtil.Errorf("WebSocket 服务启动失败: %s", err)
		}
	}()
	logUtil.Infof("WebSocket 服务已启动:%v", webSocketConfig.Addr)
}

// WebSocket 连接处理
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := wsClient.WebSocketClient.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		logUtil.Errorf("WebSocket 连接失败: %s", err)
		return
	}
	defer conn.Close()

	// 获取用户ID
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("id is required"))
		return
	}
	// 存储连接
	wsClient.WebSocketClient.Connections.Store(id, conn)
	// 连接成功后的回调
	onOpen(conn, id)
	// 处理 WebSocket 消息
	for {
		// 读取消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			onError(conn, id, err)
			break
		}
		// 消息处理函数
		wsHandler.WebSocketHandlerInstance.MessageHandler(id, msg)
	}
	// 连接断开时调用 onClose 并删除连接
	onClose(conn, id)
}

func onOpen(conn *websocket.Conn, id int64) {
	logUtil.Infof("WebSocket 客户端(%v)已连接: %s", conn.RemoteAddr(), id)
	//todo 更新心跳 时间
	wsClient.WebSocketClient.SendMessageToOne(id, "连接成功")
}

func onClose(conn *websocket.Conn, id int64) {
	logUtil.Infof("WebSocket 客户端(%v)已断开: %s", id, conn.RemoteAddr())
	wsClient.WebSocketClient.Connections.Delete(id)
}

func onError(conn *websocket.Conn, id int64, err error) {
	logUtil.Infof("WebSocket 客户端(%v)%s发生错误:%v", id, conn.RemoteAddr(), err)
}
