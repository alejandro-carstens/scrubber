package actions

import (
	"os"
	"scrubber/actions/contexts"
	"scrubber/logger"
	"scrubber/notifications"
	"strconv"

	"github.com/alejandro-carstens/golastic"
)

const DEFAULT_HEALTH_CHECK_INTERVAL int64 = 30

type action struct {
	retryCount     int
	name           string
	queue          *notifications.Queue
	context        contexts.Contextable
	reporter       *reporter
	errorReportMap *errorReportMap
	connection     *golastic.Connection
	indexer        *golastic.Indexer
}

// Init initializes an action
func (a *action) Init(context contexts.Contextable, logger *logger.Logger, connection *golastic.Connection, queue *notifications.Queue) error {
	if connection == nil {
		var healthCheckInterval int64 = DEFAULT_HEALTH_CHECK_INTERVAL

		if len(os.Getenv("ELASTICSEARCH_HEALTH_CHECK_INTERVAL")) > 0 {
			value, err := strconv.Atoi(os.Getenv("ELASTICSEARCH_HEALTH_CHECK_INTERVAL"))

			if err != nil {
				return err
			}

			healthCheckInterval = int64(value)
		}

		connection = golastic.NewConnection(&golastic.ConnectionContext{
			Urls:                []string{os.Getenv("ELASTICSEARCH_URI")},
			Password:            os.Getenv("ELASTICSEARCH_PASSWORD"),
			Username:            os.Getenv("ELASTICSEARCH_USERNAME"),
			HealthCheckInterval: healthCheckInterval,
			ErrorLogPrefix:      os.Getenv("ELASTICSEARCH_ERROR_LOG_PREFIX"),
			InfoLogPrefix:       os.Getenv("ELASTICSEARCH_INFO_LOG_PREFIX"),
		})

		if err := connection.Connect(); err != nil {
			return err
		}
	}

	a.connection = connection
	a.indexer = connection.Indexer(nil)
	a.context = context
	a.name = context.Action()
	a.errorReportMap = newErrorReportMap()
	a.reporter = newReporter(logger)
	a.queue = queue

	return nil
}

// HasErrors signals and logs whether the action experienced errors
func (a *action) HasErrors() bool {
	for _, report := range a.errorReportMap.reports {
		if a.reporter.Logger() == nil {
			break
		}

		for _, err := range report.errs {
			a.reporter.Logger().Errorf("Errors: %v", err)
		}
	}

	return a.errorReportMap.hasErrors()
}

// DisableAction indicates whether or not the action should be performed
func (a *action) DisableAction() bool {
	return a.context.Options().GetDisableAction()
}

// List returns the actionable list
func (a *action) List() []string {
	return []string{}
}

// Disconnect clears the Elasticsearch connection
func (a *action) Disconnect() {
	a.connection = nil
	a.indexer = nil

	a.Release()
}

// Release releases the notification's queue
func (a *action) Release() {
	if a.queue != nil {
		a.queue.Release()
	}
}
