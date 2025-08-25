package manager

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go-chat/configs"
	"go-chat/internal/consumer"
	"go-chat/internal/utils/logUtil"
)

type RabbitMQManager struct {
	conn   *amqp.Connection
	sendCh *amqp.Channel
	recvCh *amqp.Channel
}

var RabbitClient *RabbitMQManager

// InitRabbitMQ 初始化RabbitMQ连接和通道
func InitRabbitMQ() {
	RabbitClient = &RabbitMQManager{}
	rabbitmqConfig := configs.AppConfig.Rabbitmq
	var err error
	// 连接到RabbitMQ服务器
	RabbitClient.conn, err = amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v:%v/",
		rabbitmqConfig.Username,
		rabbitmqConfig.Password,
		rabbitmqConfig.Host,
		rabbitmqConfig.Port))
	if err != nil {
		logUtil.Errorf("rabbitmq 启动失败: %s", err)
		return
	}
	// 创建发送消息的通道
	RabbitClient.sendCh, err = RabbitClient.conn.Channel()
	// 创建接收消息的通道
	RabbitClient.recvCh, err = RabbitClient.conn.Channel()

	//启动消费者
	startConsumers()
}

// SendMessage 发送任意类型的消息到指定交换机
func (rmq *RabbitMQManager) SendMessage(exchange, routingKey interface{}, message interface{}) error {
	// 将消息体序列化为 JSON 字节数组
	body, err := json.Marshal(message)
	if err != nil {
		logUtil.Errorf("消息序列化失败: %s", err)
		return err
	}

	// 发送消息到指定交换机
	err = rmq.sendCh.Publish(
		exchange.(string),   // 交换机
		routingKey.(string), // 路由键
		false,               // 是否等待确认
		false,               // 是否强制推送
		amqp.Publishing{
			ContentType: "application/json", // 设置消息类型为 JSON
			Body:        body,               // 消息体（已序列化）
		},
	)

	// 错误检查与重试机制（示例）
	if err != nil {
		logUtil.Errorf("消息发送失败: %s", err)
		return err
	}
	logUtil.Infof("消息发送成功")
	return nil
}

func startConsumers() {
	for _, c := range configs.AppConfig.Mq {
		go runConsumer(c.Exchange, c.RoutingKey, c.Queue, consumer.HandlerMap[c.Handler])
	}
}

func runConsumer(exchange, routingKey, queueName string, handler func([]byte)) {
	err := registerConsumer(exchange, routingKey, queueName, handler)
	if err != nil {
		logUtil.Errorf("消费者(%v)注册失败: %s", exchange+"-"+routingKey+"-"+queueName, err)
	}
}

// 注册消费者并监听指定队列的消息
func registerConsumer(exchange, routingKey, queueName string, messageHandler func([]byte)) error {
	// 声明交换机
	err := RabbitClient.recvCh.ExchangeDeclare(
		exchange, // 交换机名称
		"direct", // 交换机类型（根据需求选择：direct、topic、fanout等）
		true,     // 是否持久化
		false,    // 是否自动删除
		false,    // 是否排他性
		false,    // 是否阻塞
		nil,      // 额外参数
	)
	if err != nil {
		logUtil.Errorf("交换机声明失败: %s", err)
		return err
	}

	// 声明队列
	q, err := RabbitClient.recvCh.QueueDeclare(
		queueName, // 队列名字
		true,      // 是否持久化
		false,     // 是否自动删除
		false,     // 是否排他性
		false,     // 是否阻塞
		nil,       // 额外参数
	)
	if err != nil {
		logUtil.Errorf("队列声明失败: %s", err)
		return err
	}

	// 绑定队列到交换机，使用指定的路由键
	err = RabbitClient.recvCh.QueueBind(
		q.Name,     // 队列名字
		routingKey, // 路由键
		exchange,   // 交换机
		false,      // 是否阻塞
		nil,        // 额外参数
	)
	if err != nil {
		logUtil.Errorf("队列与交换机绑定失败: %s", err)
		return err
	}
	// 消费者监听消息
	msgs, err := RabbitClient.recvCh.Consume(
		q.Name, // 队列名字
		"",     // 消费者标签（空字符串表示不指定）
		true,   // 是否自动应答
		false,  // 是否排他性
		false,  // 是否阻塞
		false,  // 是否持久化
		nil,    // 额外参数
	)
	if err != nil {
		logUtil.Errorf("消费者(%v)启动失败: %s", exchange+"-"+routingKey+"-"+queueName, err)
		return err
	}
	logUtil.Infof("消费者(%v)启动成功", exchange+"-"+routingKey+"-"+queueName)
	// 阻塞处理消息,一个消费者只负责一个队列
	for msg := range msgs {
		// 调用不同的消息处理函数
		messageHandler(msg.Body)
	}

	return nil
}

// CloseRabbitMQ 关闭RabbitMQ连接
func (rmq *RabbitMQManager) CloseRabbitMQ() {
	if rmq.sendCh != nil {
		rmq.sendCh.Close()
	}
	if rmq.recvCh != nil {
		rmq.recvCh.Close()
	}
	if rmq.conn != nil {
		rmq.conn.Close()
	}
}
