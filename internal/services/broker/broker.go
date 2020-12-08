package broker

type Logger interface {
	Fatalf(tpl string, args ...interface{})
	Errorf(tpl string, args ...interface{})
}

type Broker interface {
	Subscribe(event string, h Handler) error
	Publish(msg Message, queue string) error
}

// callback to handling messages
type Handler func(msg Message) (bool, error)

// AMQP message
type Message struct {
	Body    []byte
	Headers map[string]interface{}
}
