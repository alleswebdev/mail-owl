package rabbitmq

import (
	"fmt"
	"github.com/alleswebdev/mail-owl/internal/services/broker"
	"github.com/streadway/amqp"
)

type Rabbit struct {
	Conn      *amqp.Connection
	exchange  string
	consumers []Consumer
	defaultCh *amqp.Channel
	log       broker.Logger
}

func NewRabbit(uri string, exchange string, log broker.Logger) (*Rabbit, error) {
	conn, err := amqp.Dial(uri)

	if err != nil {
		return nil, fmt.Errorf("dial: %s", err)
	}

	defaultCh, err := conn.Channel()

	if err != nil {
		return nil, fmt.Errorf("open channel: %s", err)
	}

	go func() {
		log.Fatalf("Mq connection was closed: %s", <-conn.NotifyClose(make(chan *amqp.Error)))
	}()

	return &Rabbit{Conn: conn, exchange: exchange, defaultCh: defaultCh, log: log}, nil
}

// Die if the connect is lost (k8s reloaded this pod)
func (c *Rabbit) NotifyClose(log broker.Logger) {
	go func() {
		log.Fatalf("Mq connection was closed: %s", <-c.Conn.NotifyClose(make(chan *amqp.Error)))
	}()
}

func (c *Rabbit) Subscribe(event string, h broker.Handler) error {
	ch, err := c.Conn.Channel()

	if nil != err {
		return err
	}

	cs := Consumer{Ch: ch, exchange: c.exchange}

	cs.bindingKey = event

	// 5 msg per query
	err = ch.Qos(5, 0, true)

	if nil != err {
		return err
	}

	err, queue := AnnounceQueue(ch, cs.bindingKey, cs.exchange)

	if nil != err {
		return err
	}

	cs.Queue = *queue

	events, err := cs.Ch.Consume(
		cs.Queue.Name, // name
		"",            // consumerTag,
		false,         // noAck
		false,         // exclusive
		false,         // noLocal
		false,         // noWait
		nil,           // arguments
	)

	if err != nil {
		return fmt.Errorf("queue Consume: %s", err)
	}

	c.consumers = append(c.consumers, cs)

	go func() {
		for e := range events {
			isNack, err := h(broker.Message{
				Body:    e.Body,
				Headers: e.Headers,
			})

			if err != nil {
				c.log.Errorf("msg handling error:%s", err)

				if isNack {
					err = e.Nack(false, true)
					if err != nil {
						c.log.Fatalf("msg nack error:%s", err)
					}

					continue
				}
			}

			err = e.Ack(false)

			if err != nil {
				c.log.Fatalf("msg ack error:%s", err)
			}
		}

		defer ch.Close()
	}()
	return nil
}

func (c *Rabbit) Publish(msg broker.Message, queue string) error {
	var err error

	if c.defaultCh == nil {
		c.defaultCh, err = c.Conn.Channel()
	}

	if err != nil {
		return err
	}

	err, _ = AnnounceQueue(c.defaultCh, queue, c.exchange)

	if err != nil {
		return err
	}

	err = c.defaultCh.Publish(
		c.exchange,
		queue, // routing key
		false, // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         msg.Body,
			Headers:      msg.Headers,
		})

	if err != nil {
		return err
	}

	return nil
}
