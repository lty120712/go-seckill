package wsMessage

import "time"

type Message struct {
	Type   string      `json:"type"`    // 事件类型
	SendId int64       `json:"send_id"` // 发送者ID
	Data   interface{} `json:"data"`    // 具体数据
	Time   time.Time   `json:"time"`    //  消息发送时间
}

// 事件类型
const (
	Chat         = "chat"          //聊天
	ChatAck      = "chat_ack"      // 聊天确认
	OnlineStatus = "online_status" // 在线状态
	Recall       = "recall"        //  撤回
	IdRequest    = "id_request"    // 请求获取真实ID,引入mq之后采用

	HeartBeat = "heartbeat" //心跳检测

	HeartBeatAck = "heartbeat_ack" //心跳检测确认
)
