package configurations

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Slack struct {
	Webhooks   map[string]string
	RetryCount int
}

func (s *Slack) FillFromEnvs() Configurable {
	webhooks := map[string]string{}

	for _, webhook := range strings.Split(os.Getenv("SLACK_WEBHOOKS"), ",") {
		parts := strings.Split(webhook, "@")

		if len(parts) != 2 {
			continue
		}

		webhooks[parts[0]] = parts[1]
	}

	s.Webhooks = webhooks

	retryCount, _ := strconv.Atoi(os.Getenv("SLACK_RETRY_COUNT"))

	s.RetryCount = retryCount

	return s
}

func (s *Slack) Validate() (Configurable, error) {
	if len(s.Webhooks) == 0 {
		return nil, errors.New("You must specify at least one account")
	}

	for _, webhook := range s.Webhooks {
		if len(webhook) == 0 {
			return nil, errors.New("you must specify a webhook per account")
		}
	}

	return s, nil
}
