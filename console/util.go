package console

import (
	"context"

	"scrubber/actions"
	"scrubber/actions/contexts"
	"scrubber/logger"
	"scrubber/notifications"
	rp "scrubber/resourcepool"

	"github.com/alejandro-carstens/golastic"
)

// NewScheduler instantiates a scheduler
func NewScheduler(
	basePath string,
	exclude []string,
	queue *notifications.Queue,
) *scheduler {
	return &scheduler{
		basePath:   basePath,
		exclude:    exclude,
		logger:     rp.Logger(),
		connection: rp.Elasticsearch(),
		queue:      queue,
		context:    rp.Context(),
	}
}

// Execute performs a given action
func Execute(
	context contexts.Contextable,
	logger *logger.Logger,
	connection *golastic.Connection,
	queue *notifications.Queue,
	ctx context.Context,
) {
	action, err := actions.Create(context, logger, connection, queue, ctx)

	if err != nil {
		logger.Errorf(err.Error())

		return
	}

	defer action.Disconnect()

	if action.DisableAction() {
		logger.Noticef("%v action disabled", context.Action())

		return
	}

	if !action.Perform().HasErrors() {
		logger.Noticef("successfully executed %v action", context.Action())
	}

	if err := action.Notify(); err != nil {
		logger.Errorf("an error [%v] occurred while notifying action", err.Error())
	}
}
