package notifications

import (
	"scrubber/notifications/configurations"
	"scrubber/notifications/messages"
)

type Notifiable interface {
	Configure(configuration configurations.Configurable) error

	Send(message messages.Sendable) error

	Retry() error
}
