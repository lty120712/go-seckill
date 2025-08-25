package interfaces

type RabbitMQManager interface {
	SendMessage(exchange, routingKey interface{}, message interface{}) error
}
