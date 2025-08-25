package consumer

// 注册消费者处理函数
var HandlerMap = map[string]func([]byte){
	"HandleOrder": HandleOrderConsumer,
}

// HandleOrderConsumer 处理 order 队列的消息
func HandleOrderConsumer(msg []byte) {
	return
}
