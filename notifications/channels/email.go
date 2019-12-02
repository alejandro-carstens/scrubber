package channels

import (
	"errors"
	"log"
	"scrubber/notifications/configurations"
	"scrubber/notifications/messages"

	"github.com/cenkalti/backoff"
	"github.com/go-mail/mail"
)

// Email represents an e-mail notification channel
type Email struct {
	configuration *configurations.Email
	message       *messages.Email
	retryCount    int
}

// Configure is responsible for configuring the notification channel
func (e *Email) Configure(configuration configurations.Configurable) error {
	config, valid := configuration.(*configurations.Email)

	if !valid {
		return errors.New("Invalid configuration, not of type email")
	}

	e.configuration = config

	return nil
}

// Send is redsponsible for sending the notification over the selected channel
func (e *Email) Send(message messages.Sendable) error {
	msg, valid := message.(*messages.Email)

	if !valid {
		return errors.New("Invalid configuration, not of type email")
	}

	e.message = msg

	dialer := mail.NewDialer(
		e.configuration.Host,
		e.configuration.Port,
		e.configuration.Username,
		e.configuration.Password,
	)

	dialer.StartTLSPolicy = mail.MandatoryStartTLS

	sender, err := dialer.Dial()

	if err != nil {
		return err
	}

	m := mail.NewMessage()

	m.SetHeader("From", msg.From)
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html", msg.Body)

	err = mail.Send(sender, m)

	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("No errors")
	}

	return err
}

// Retry is responsible for trying to complete the notification in case errors occur
func (e *Email) Retry() error {
	if e.message == nil {
		return errors.New("message not set, please verify the webhook configuration")
	}

	return backoff.Retry(func() error {
		if err := e.Send(e.message); err != nil {
			e.retryCount++

			if e.retryCount == e.configuration.RetryCount {
				return &backoff.PermanentError{
					Err: err,
				}
			}

			return err
		}

		return nil
	}, backoff.NewExponentialBackOff())
}
