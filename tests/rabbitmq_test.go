// rabbitmq_test.go
package tests

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"testing"
)

// 测试是否能够连接到 RabbitMQ 并进行基本操作
func TestRabbitMQ_Connection(t *testing.T) {
	conn, err := amqp.Dial("amqp://lty:lty0712@8.137.38.55:5672/")
	if err != nil {
		t.Fatalf("无法连接到 RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("无法创建通道: %v", err)
	}
	defer ch.Close()

	// 声明一个队列
	q, err := ch.QueueDeclare(
		"test_queue", // 队列名称
		true,         // 是否持久化
		false,        // 是否自动删除
		false,        // 是否排他性
		false,        // 是否阻塞
		nil,          // 额外参数
	)
	if err != nil {
		t.Fatalf("无法声明队列: %v", err)
	}

	// 验证队列是否存在
	if q.Name != "test_queue" {
		t.Fatalf("队列名称不正确: %v", q.Name)
	}

	t.Log("RabbitMQ 连接和队列声明成功")
}
