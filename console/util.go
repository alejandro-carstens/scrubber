package console

import (
	"scrubber/actions"
	"scrubber/actions/contexts"
	"scrubber/logger"
)

func NewScheduler(basePath string, logger *logger.Logger) *Scheduler {
	return &Scheduler{basePath: basePath, logger: logger}
}

func Execute(context contexts.Contextable, logger *logger.Logger) {
	action, err := actions.Create(context, logger)

	if err != nil {
		logger.Errorf("%v", err.Error())
		return
	}

	defer action.TearDownBuilder()

	if action.DisableAction() {
		logger.Noticef("%v action disabled", context.Action())
		return
	}

	if !action.Perform().HasErrors() {
		logger.Noticef("successfully executed %v action", context.Action())
	}
}
