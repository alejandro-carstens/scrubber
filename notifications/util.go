package notifications

import (
	"github.com/alejandro-carstens/scrubber/logger"
	"github.com/alejandro-carstens/scrubber/notifications/channels"
	"github.com/alejandro-carstens/scrubber/notifications/messages"
)

// NewQueue instantiates a queue
func NewQueue(capacity int, logger *logger.Logger) (*Queue, error) {
	return new(Queue).Init(capacity, logger)
}

// Notify sends a notification over a Notifiable channel
func Notify(message messages.Sendable) error {
	return channels.Notify(message)
}
