package notifications

import (
	"scrubber/logger"
	"scrubber/notifications/channels"
	"scrubber/notifications/messages"
)

// NewQueue instantiates a queue
func NewQueue(capacity int, logger *logger.Logger) (*Queue, error) {
	return new(Queue).Init(capacity, logger)
}

// Notify sends a notification over a Notifiable channel
func Notify(message messages.Sendable) error {
	return channels.Notify(message)
}
