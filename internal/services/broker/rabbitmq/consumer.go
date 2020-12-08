package rabbitmq

import (
	"github.com/streadway/amqp"
)

type Consumer struct {
	Conn       amqp.Connection
	Ch         *amqp.Channel
	Queue      amqp.Queue
	exchange   string // exchange that we will bind to
	bindingKey string // routing key that we are using
}
