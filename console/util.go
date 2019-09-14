package console

import (
	"scrubber/actions"
	"scrubber/actions/contexts"
	"scrubber/logger"
	"scrubber/notifications/queue"

	"github.com/alejandro-carstens/golastic"
)

func NewScheduler(basePath string, logger *logger.Logger, builder *golastic.Connection, enqueuer *queue.Enqueuer) *Scheduler {
	return &Scheduler{basePath: basePath, logger: logger, builder: builder, enqueuer: enqueuer}
}

func Execute(context contexts.Contextable, logger *logger.Logger, builder *golastic.Connection, enqueuer *queue.Enqueuer) {
	action, err := actions.Create(context, logger, builder, enqueuer)

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
