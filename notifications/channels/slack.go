package channels

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"scrubber/notifications/configurations"
	"scrubber/notifications/messages"
	"strconv"
	"time"

	"github.com/cenkalti/backoff"
)

type webhookMessage struct {
	Attachments []attachment `json:"attachments,omitempty"`
}

type attachment struct {
	Color         string            `json:"color,omitempty"`
	Fallback      string            `json:"fallback"`
	AuthorName    string            `json:"author_name,omitempty"`
	AuthorSubname string            `json:"author_subname,omitempty"`
	AuthorIcon    string            `json:"author_icon,omitempty"`
	Text          string            `json:"text"`
	Fields        []attachmentField `json:"fields,omitempty"`
	Footer        string            `json:"footer,omitempty"`
	FooterIcon    string            `json:"footer_icon,omitempty"`
	Ts            json.Number       `json:"ts,omitempty"`
}

type attachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// Slack represents a Slack notification channel
type Slack struct {
	configuration *configurations.Slack
	message       *messages.Slack
	retryCount    int
}

// Configure is responsible for configuring the notification channel
func (s *Slack) Configure(configuration configurations.Configurable) error {
	config, valid := configuration.(*configurations.Slack)

	if !valid {
		return errors.New("Invalid configuration, not of type slack")
	}

	s.configuration = config

	return nil
}

// Send is redsponsible for sending the notification over the selected channel
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

	return s.postWebhook(webhook, &webhookMessage{
		Attachments: []attachment{
			attachment{
				Color:         msg.Attachment.Color,
				Fallback:      msg.Attachment.Fallback,
				AuthorName:    msg.Attachment.AuthorName,
				AuthorSubname: msg.Attachment.AuthorSubname,
				AuthorIcon:    msg.Attachment.AuthorIcon,
				Text:          msg.Attachment.Text,
				Footer:        msg.Attachment.Footer,
				FooterIcon:    msg.Attachment.FooterIcon,
				Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
				Fields: []attachmentField{
					attachmentField{
						Title: "Message Key",
						Value: msg.DedupKey,
					},
				},
			},
		},
	})
}

// Retry is responsible for trying to complete the notification in case errors occur
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

func (s *Slack) postWebhook(url string, msg *webhookMessage) error {
	raw, err := json.Marshal(msg)

	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Post(url, "application/json", bytes.NewReader(raw))

	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New("Received a non 200 status code")
	}

	return nil
}
