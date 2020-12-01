package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

type Consumer struct {
	Conn       amqp.Connection
	Ch         amqp.Channel
	Queue      amqp.Queue
	exchange   string // exchange that we will bind to
	bindingKey string // routing key that we are using
}

// declare queue
func (c *Consumer) AnnounceQueue() error {
	if err := c.Ch.ExchangeDeclare(
		c.exchange, // name of the exchange
		"direct",   // type
		true,       // durable
		false,      // delete when complete
		false,      // internal
		false,      // noWait
		nil,
	); err != nil {
		return fmt.Errorf("exchange Declare: %s", err)
	}

	queue, err := c.Ch.QueueDeclare(
		c.bindingKey, // name of the queue
		true,         // durable
		false,        // delete when usused
		false,        // exclusive
		false,        // noWait
		nil,          // arguments
	)

	if err != nil {
		return fmt.Errorf("queue declare: %s", err)
	}

	c.Queue = queue

	if err = c.Ch.QueueBind(
		queue.Name, // name of the queue
		queue.Name, // bindingKey
		c.exchange, // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return fmt.Errorf("queue bind: %s", err)
	}

	return nil
}
