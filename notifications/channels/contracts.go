package channels

import (
	"github.com/alejandro-carstens/scrubber/notifications/configurations"
	"github.com/alejandro-carstens/scrubber/notifications/messages"
)

// Notifiable represents the contract that
// notfication channels should comply with
type Notifiable interface {
	// Configure is responsible for configuring the notification channel
	Configure(configuration configurations.Configurable) error

	// Send is redsponsible for sending the notification over the selected channel
	Send(message messages.Sendable) error

	// Retry is responsible for trying to complete the notification in case errors occur
	Retry() error
}
