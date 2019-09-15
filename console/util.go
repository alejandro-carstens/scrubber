package console

import (
	"scrubber/actions"
	"scrubber/actions/contexts"
	"scrubber/logger"
	"scrubber/notifications"

	"github.com/alejandro-carstens/golastic"
)

func NewScheduler(basePath string, logger *logger.Logger, builder *golastic.Connection, queue *notifications.Queue) *Scheduler {
	return &Scheduler{basePath: basePath, logger: logger, builder: builder, queue: queue}
}

func Execute(context contexts.Contextable, logger *logger.Logger, builder *golastic.Connection, queue *notifications.Queue) {
	action, err := actions.Create(context, logger, builder, queue)

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
}
