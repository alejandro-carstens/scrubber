package configurations

import (
	"errors"
	"os"
	"strconv"
)

// Email represents the configuration required
// for sending a notification over e-mail
type Email struct {
	RetryCount int
	Port       int
	Host       string
	Username   string
	Password   string
}

// FillFromEnvs is responsible for setting the configuration
// for the channel from the respective env variables
func (e *Email) FillFromEnvs() Configurable {
	retryCount, err := strconv.Atoi(os.Getenv("EMAIL_RETRY_COUNT"))

	if err == nil {
		e.RetryCount = retryCount
	}

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))

	if err == nil {
		e.Port = port
	}

	e.Host = os.Getenv("SMTP_HOST")
	e.Username = os.Getenv("EMAIL_USERNAME")
	e.Password = os.Getenv("EMAIL_PASSWORD")

	return e
}

// Validate validates the configuration for a given channel
func (e *Email) Validate() (Configurable, error) {
	if len(e.Host) == 0 {
		return nil, errors.New("the e-mail smpt_host needs to be set")
	}
	if e.Port == 0 {
		return nil, errors.New("the e-mail smpt_port needs to be set")
	}
	if len(e.Username) == 0 {
		return nil, errors.New("the e-mail username needs to be set")
	}
	if len(e.Password) == 0 {
		return nil, errors.New("the e-mail password needs to be set")
	}

	return e, nil
}
