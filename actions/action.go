package actions

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/contexts"
	"github.com/alejandro-carstens/scrubber/logger"
	"github.com/alejandro-carstens/scrubber/notifications"
	"github.com/alejandro-carstens/scrubber/notifications/messages"
)

const DEFAULT_HEALTH_CHECK_INTERVAL int64 = 30

type action struct {
	retryCount     int
	name           string
	notifiableList []string
	queue          *notifications.Queue
	context        contexts.Contextable
	reporter       *reporter
	errorContainer *errorContainer
	connection     *golastic.Connection
	indexer        *golastic.Indexer
}

// Init initializes an action
func (a *action) Init(
	ctx contexts.Contextable,
	logger *logger.Logger,
	connection *golastic.Connection,
	queue *notifications.Queue,
) error {
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
			Context:             context.Background(),
		})

		if err := connection.Connect(); err != nil {
			return err
		}
	}

	a.connection = connection
	a.indexer = connection.Indexer(nil)
	a.context = ctx
	a.name = ctx.Action()
	a.errorContainer = newErrorContainer()
	a.reporter = newReporter(logger)
	a.queue = queue
	a.notifiableList = []string{}

	return nil
}

// HasErrors signals and logs whether the action experienced errors
func (a *action) HasErrors() bool {
	for _, report := range a.errorContainer.reports {
		if a.reporter.logger == nil {
			break
		}

		for _, err := range report.errors {
			a.reporter.logger.Errorf("Errors: %v", err)
		}
	}

	return a.errorContainer.hasErrors()
}

// DisableAction indicates whether or not the action should be performed
func (a *action) DisableAction() bool {
	return a.context.Options().GetDisableAction()
}

// Notify issues a notification regarding the execution
// of an action over an actionable list
func (a *action) Notify() error {
	if len(a.notifiableList) == 0 {
		return nil
	}

	if !a.context.Options().IsNotifiable() {
		return nil
	}

	notification := &actionNotification{
		Text: fmt.Sprintf("Successfully executed %v action for: %v", a.name, strings.Join(a.notifiableList, ", ")),
	}

	if err := json.Unmarshal(a.context.Options().GetContainer().Bytes(), notification); err != nil {
		return err
	}

	h := sha256.New()
	h.Write([]byte(notification.payload().String()))

	message, err := messages.NewMessage(notification.payload(), nil, fmt.Sprintf("%x", h.Sum(nil)))

	if err != nil {
		return err
	}

	return a.queue.Push(message)
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
