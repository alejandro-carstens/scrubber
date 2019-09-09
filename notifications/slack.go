package notifications

import (
	"encoding/json"
	"errors"
	"scrubber/notifications/configurations"
	"scrubber/notifications/messages"
	"strconv"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/nlopes/slack"
)

type Slack struct {
	configuration *configurations.Slack
	message       *messages.Slack
	retryCount    int
}

func (s *Slack) Configure(configuration configurations.Configurable) error {
	config, valid := configuration.(*configurations.Slack)

	if !valid {
		return errors.New("Invalid configuration, not of type slack")
	}

	s.configuration = config

	return nil
}

func (s *Slack) Send(message messages.Sendable) error {
	msg, valid := message.(*messages.Slack)

	if !valid {
		return errors.New("invalid message type, not a slack message")
	}

	webhook, valid := s.configuration.Webhooks[msg.Attachment.WebhookName]

	if !valid {
		return errors.New("could not find the account to send message")
	}

	s.message = msg

	attachment := slack.Attachment{
		Color:         msg.Attachment.Color,
		Fallback:      msg.Attachment.Fallback,
		AuthorName:    msg.Attachment.AuthorName,
		AuthorSubname: msg.Attachment.AuthorSubname,
		AuthorIcon:    msg.Attachment.AuthorIcon,
		Text:          msg.Attachment.Text,
		Footer:        msg.Attachment.Footer,
		FooterIcon:    msg.Attachment.FooterIcon,
		Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}

	return slack.PostWebhook(webhook, &slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	})
}

func (s *Slack) Retry() error {
	if s.message == nil {
		return errors.New("message not set, please verify the webhook configuration")
	}

	return backoff.Retry(func() error {
		if err := s.Send(s.message); err != nil {
			s.retryCount++

			if s.retryCount == s.configuration.RetryCount {
				return &backoff.PermanentError{
					Err: err,
				}
			}

			return err
		}

		return nil
	}, backoff.NewExponentialBackOff())
}
