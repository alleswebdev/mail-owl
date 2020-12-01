package broker

type Logger interface {
	Fatalf(tpl string, args ...interface{})
}

type Broker interface {
	Subscribe(event string, h Handler) error
	Publish(event string, msg Message) error
}

// callback to handling messages
type Handler func(msg Message) error

// AMQP message
type Message struct {
	Body    []byte
	Headers map[string]interface{}
}
