package messages

import (
	"errors"

	"github.com/Jeffail/gabs"
)

// NewMessage returns a sendable message
func NewMessage(payload *gabs.Container, context interface{}, dedupKey string) (Sendable, error) {
	notificationChannel, valid := payload.S("notification_channel").Data().(string)

	if !valid {
		return nil, errors.New("invalid message payload")
	}

	var message Sendable

	switch notificationChannel {
	case "slack":
		message = &Slack{}
		break
	case "pager_duty":
		message = &PagerDuty{}
		break
	case "email":
		message = &Email{}
		break
	default:
		return nil, errors.New("invalid message type")
	}

	message.Init(context, dedupKey, payload)

	if err := message.Format(); err != nil {
		return nil, err
	}

	return message, nil
}
