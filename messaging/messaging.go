package messaging

// Messaging represents the messaging system interface
type Messaging interface {
	Publish(topic string, data []byte) error
	Subscribe(topic string, handler func(data []byte)) error
	Close() error
}
