package channels

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/alejandro-carstens/scrubber/notifications/configurations"
	"github.com/alejandro-carstens/scrubber/notifications/messages"
	"github.com/cenkalti/backoff"
)

// PAGER_DUTY_BASE_URI represents the uri
// for which we need to issue the post
const PAGER_DUTY_BASE_URI string = "https://events.pagerduty.com/v2/enqueue"

// PagerDuty represents a Pager Duty notification channel
type PagerDuty struct {
	configuration *configurations.PagerDuty
	message       *messages.PagerDuty
	retryCount    int
}

// Configure is responsible for configuring the notification channel
func (pd *PagerDuty) Configure(configuration configurations.Configurable) error {
	config, valid := configuration.(*configurations.PagerDuty)

	if !valid {
		return errors.New("Invalid configuration, not of type pager_duty")
	}

	pd.configuration = config

	return nil
}

// Send is redsponsible for sending the notification over the selected channel
func (pd *PagerDuty) Send(message messages.Sendable) error {
	msg, valid := message.(*messages.PagerDuty)

	if !valid {
		return errors.New("invalid message type, not a pager_duty message")
	}

	pd.message = msg

	raw, err := json.Marshal(msg.Event)

	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", pd.getBaseUri(), bytes.NewReader(raw))

	if err != nil {
		return err
	}

	request.Header.Add("X-Routing-Key", pd.configuration.RoutingKey)

	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New("Received a non 200 status code")
	}

	return nil
}

// Retry is responsible for trying to complete the notification in case errors occur
func (pd *PagerDuty) Retry() error {
	if pd.message == nil {
		return errors.New("message not set, please verify the webhook configuration")
	}

	return backoff.Retry(func() error {
		if err := pd.Send(pd.message); err != nil {
			pd.retryCount++

			if pd.retryCount == pd.configuration.RetryCount {
				return &backoff.PermanentError{
					Err: err,
				}
			}

			return err
		}

		return nil
	}, backoff.NewExponentialBackOff())
}

func (pd *PagerDuty) getBaseUri() string {
	if len(os.Getenv("PAGER_DUTY_BASE_URI")) == 0 {
		return PAGER_DUTY_BASE_URI
	}

	return os.Getenv("PAGER_DUTY_BASE_URI")
}
