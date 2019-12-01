package configurations

import (
	"errors"
	"os"
	"strconv"
)

// PagerDuty represents the configuration required
// for sending a message over to PagerDuty
type PagerDuty struct {
	RoutingKey string
	RetryCount int
}

// FillFromEnvs is responsible for setting the configuration
// for the channel from the respective env variables
func (pd *PagerDuty) FillFromEnvs() Configurable {
	pd.RoutingKey = os.Getenv("PAGER_DUTY_ROUTING_KEY")

	retryCount, _ := strconv.Atoi(os.Getenv("PAGER_DUTY_RETRY_COUNT"))

	pd.RetryCount = retryCount

	return pd
}

// Validate validates the configuration for a given channel
func (pd *PagerDuty) Validate() (Configurable, error) {
	if len(pd.RoutingKey) == 0 {
		return nil, errors.New("the Pager Duty routing_key needs to be set")
	}

	return pd, nil
}
