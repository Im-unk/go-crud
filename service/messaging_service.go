package service

import (
	"main.go/messaging"
)

// MessagingService handles the messaging functionality
type MessagingService struct {
	messaging messaging.Messaging
}

// NewMessagingService creates a new MessagingService
func NewMessagingService(messaging messaging.Messaging) *MessagingService {
	return &MessagingService{
		messaging: messaging,
	}
}

// Publish publishes a message using the messaging system
func (s *MessagingService) Publish(topic string, data []byte) error {
	return s.messaging.Publish(topic, data)
}

// Subscribe subscribes to a topic and registers a message handler
func (s *MessagingService) Subscribe(topic string, handler func(data []byte)) error {
	return s.messaging.Subscribe(topic, handler)
}

// Close closes the messaging connection
func (s *MessagingService) Close() error {
	return s.messaging.Close()
}
