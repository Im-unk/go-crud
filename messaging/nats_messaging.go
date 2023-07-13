package messaging

import (
	"github.com/nats-io/nats.go"
)

// NatsMessaging represents the NATS messaging system
type NatsMessaging struct {
	nc *nats.Conn
}

// NewNatsMessaging creates a new instance of NatsMessaging
func NewNatsMessaging(url string) (Messaging, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NatsMessaging{nc: nc}, nil
}

// Publish publishes a message to a NATS topic
func (n *NatsMessaging) Publish(topic string, data []byte) error {
	return n.nc.Publish(topic, data)
}

// Subscribe subscribes to a NATS topic and registers a message handler
func (n *NatsMessaging) Subscribe(topic string, handler func(data []byte)) error {
	_, err := n.nc.Subscribe(topic, func(msg *nats.Msg) {
		handler(msg.Data)
	})
	return err
}

// Close closes the NATS connection
func (n *NatsMessaging) Close() error {
	n.nc.Close()
	return nil
}
