package channels

import (
	"github.com/alejandro-carstens/scrubber/notifications/configurations"
	"github.com/alejandro-carstens/scrubber/notifications/messages"
)

// Notify sends a notification over a Notifiable channel
func Notify(message messages.Sendable) error {
	var channel Notifiable

	switch message.Type() {
	case "slack":
		channel = &Slack{}
		break
	case "pager_duty":
		channel = &PagerDuty{}
		break
	case "email":
		channel = &Email{}
		break
	}

	config, err := configurations.Config(message.Type())

	if err != nil {
		return err
	}

	if err := channel.Configure(config); err != nil {
		return err
	}

	if err := channel.Send(message); err != nil {
		return channel.Retry()
	}

	return nil
}
