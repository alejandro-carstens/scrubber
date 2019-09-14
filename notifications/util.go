package notifications

import (
	"scrubber/logger"
	"scrubber/notifications/channels"
	"scrubber/notifications/messages"
	"scrubber/notifications/queue"
)

// NewEnqueuer instantiates an enqueuer
func NewEnqueuer(capacity int, logger *logger.Logger) (*queue.Enqueuer, error) {
	return new(queue.Enqueuer).Init(capacity, logger)
}

// Notify sends a notification over a Notifiable channel
func Notify(message messages.Sendable) error {
	return channels.Notify(message)
}
