package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

const DEFAULT_HEALTH_CHECK_INTERVAL int64 = 30

// Elasticsearch represents an Elastisearch
// connection configuration
type Elasticsearch struct {
	Urls                []string
	Sniff               bool
	HealthCheckInterval int64
	ErrorLogPrefix      string
	InfoLogPrefix       string
	Password            string
	Username            string
}

// FillFromEnvs fills the Elasticsearch configuration
// from environmental variables
func (e *Elasticsearch) FillFromEnvs() Configurable {
	e.Urls = strings.Split(os.Getenv("ELASTICSEARCH_URI"), ",")
	e.ErrorLogPrefix = os.Getenv("ELASTICSEARCH_ERROR_LOG_PREFIX")
	e.InfoLogPrefix = os.Getenv("ELASTICSEARCH_INFO_LOG_PREFIX")
	e.Password = os.Getenv("ELASTIC_PASSWORD")
	e.Username = os.Getenv("ELASTIC_USERNAME")

	healthCheckInterval, err := strconv.Atoi(os.Getenv("ELASTICSEARCH_HEALTH_CHECK_INTERVAL"))

	if err != nil {
		e.HealthCheckInterval = DEFAULT_HEALTH_CHECK_INTERVAL
	} else {
		e.HealthCheckInterval = int64(healthCheckInterval)
	}

	sniff, err := strconv.ParseBool(os.Getenv("ELASTICSEARCH_SNIFF"))

	if err == nil {
		e.Sniff = sniff
	}

	return e
}

// Validate validates the Elasticsearch configuration
func (e *Elasticsearch) Validate() (Configurable, error) {
	if len(e.Urls) == 0 {
		return nil, errors.New("at least one url is required in order to connect to elasticsearch")
	}

	if e.HealthCheckInterval <= 0 {
		return nil, errors.New("healthcheck interval needs to be greater than 0")
	}

	return e, nil
}
