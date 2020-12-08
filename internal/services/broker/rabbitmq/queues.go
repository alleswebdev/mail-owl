package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

const SchedulerQueue = "mailowl-scheduler"
const BuilderQueue = "mailowl-builder"

// declare queue
func AnnounceQueue(ch *amqp.Channel, name string, exchange string) (error, *amqp.Queue) {
	if err := ch.ExchangeDeclare(
		exchange, // name of the exchange
		"direct", // type
		true,     // durable
		false,    // delete when complete
		false,    // internal
		false,    // noWait
		nil,
	); err != nil {
		return fmt.Errorf("exchange Declare: %s", err), nil
	}

	queue, err := ch.QueueDeclare(
		name,  // name of the queue
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)

	if err != nil {
		return fmt.Errorf("queue declare: %s", err), nil
	}

	if err = ch.QueueBind(
		queue.Name, // name of the queue
		queue.Name, // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return fmt.Errorf("queue bind: %s", err), nil
	}

	return nil, &queue
}
