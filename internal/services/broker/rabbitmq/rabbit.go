package rabbitmq

import (
	"fmt"
	"github.com/alleswebdev/mail-owl/internal/services/broker"
	"github.com/sirupsen/logrus"
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

	return &Rabbit{Conn: conn, exchange: exchange, defaultCh: defaultCh, log: log}, nil
}

// Die if the connect is lost (k8s reloaded this pod)
func (c *Rabbit) NotifyClose(log broker.Logger) {
	go func() {
		log.Fatalf("Mq connection was closed: %s", <-c.Conn.NotifyClose(make(chan *amqp.Error)))
	}()
}

func (c *Rabbit) Subscribe(event string, h broker.Handler) (*Consumer, error) {
	ch, err := c.Conn.Channel()

	if nil != err {
		return nil, err
	}

	c.NotifyClose(c.log)

	cs := Consumer{Ch: *ch, exchange: c.exchange}

	cs.bindingKey = event

	// 5 msg per query
	err = ch.Qos(5, 0, true)

	if nil != err {
		return nil, err
	}

	err = cs.AnnounceQueue()

	if nil != err {
		return nil, err
	}

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
		return nil, fmt.Errorf("queue Consume: %s", err)
	}

	c.consumers = append(c.consumers, cs)

	go func() {
		for e := range events {
			err = h(broker.Message{
				Body:    e.Body,
				Headers: e.Headers,
			})

			if err != nil {
				logrus.Errorf("msg handling error:%s", err)
				err = e.Nack(false, true)
			}

			err = e.Ack(false)

			if err != nil {
				logrus.Fatalf("msg ack/nack error:%s", err)
			}
		}

		defer ch.Close()
	}()
	return &cs, nil
}

func ProducerHandler(msg broker.Message) error {
	fmt.Println(msg.Body)
	fmt.Println(msg.Headers)
	fmt.Println("Finish!")

	return nil
}

func (c *Rabbit) Publish(msg broker.Message, queue string) error {
	var err error

	if c.defaultCh == nil {
		c.defaultCh, err = c.Conn.Channel()
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
