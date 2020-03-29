package console

import (
	"context"

	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions"
	"github.com/alejandro-carstens/scrubber/actions/contexts"
	"github.com/alejandro-carstens/scrubber/logger"
	"github.com/alejandro-carstens/scrubber/notifications"
)

// NewScheduler instantiates a scheduler
func NewScheduler(
	basePath string,
	exclude []string,
	logger *logger.Logger,
	builder *golastic.Connection,
	queue *notifications.Queue,
	context context.Context,
) *scheduler {
	return &scheduler{
		basePath: basePath,
		exclude:  exclude,
		logger:   logger,
		builder:  builder,
		queue:    queue,
		context:  context,
	}
}

// Execute performs a given action
func Execute(
	context contexts.Contextable,
	logger *logger.Logger,
	builder *golastic.Connection,
	queue *notifications.Queue,
	ctx context.Context,
) {
	action, err := actions.Create(context, logger, builder, queue, ctx)

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
