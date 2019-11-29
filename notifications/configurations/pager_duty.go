package configurations

import (
	"errors"
	"os"
	"strconv"
)

// PagerDuty represents the configuration required
// for sending a message over to PagerDuty
type PagerDuty struct {
	ApiKey     string
	RetryCount int
}

// FillFromEnvs is responsible for setting the configuration
// for the channel from the respective env variables
func (pd *PagerDuty) FillFromEnvs() Configurable {
	pd.ApiKey = os.Getenv("PAGER_DUTY_API_KEY")

	retryCount, _ := strconv.Atoi(os.Getenv("PAGER_DUTY_RETRY_COUNT"))

	pd.RetryCount = retryCount

	return pd
}

// Validate validates the configuration for a given channel
func (pd *PagerDuty) Validate() (Configurable, error) {
	if len(pd.ApiKey) == 0 {
		return nil, errors.New("the Pager Duty api_key needs to be set")
	}

	return pd, nil
}
