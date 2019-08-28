package console

import (
	"scrubber/actions"
	"scrubber/actions/contexts"
	"scrubber/logger"

	"github.com/alejandro-carstens/golastic"
)

func NewScheduler(basePath string, logger *logger.Logger, builder *golastic.Connection) *Scheduler {
	return &Scheduler{basePath: basePath, logger: logger, builder: builder}
}

func Execute(context contexts.Contextable, logger *logger.Logger, builder *golastic.Connection) {
	action, err := actions.Create(context, logger, builder)

	if err != nil {
		logger.Errorf("%v", err.Error())
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
