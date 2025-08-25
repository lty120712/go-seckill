package consumer

// 注册消费者处理函数
var HandlerMap = map[string]func([]byte){
	"HandleChat": HandleChatConsumer,
}
