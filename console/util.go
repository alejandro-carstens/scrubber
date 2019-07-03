package console

import (
	"scrubber/actions"
	"scrubber/actions/contexts"
	"scrubber/logging"
)

func NewScheduler(basePath string, logger *logging.SrvLogger) *Scheduler {
	return &Scheduler{basePath: basePath, logger: logger}
}

func Execute(context contexts.Contextable, logger *logging.SrvLogger) {
	action, err := actions.Create(context, logger)

	if err != nil {
		logger.Errorf("%v", err.Error())

		return
	}

	if action.DisableAction() {
		logger.Noticef("%v action disabled", context.Action())

		return
	}

	if !action.Perform().HasErrors() {
		logger.Noticef("successfully executed %v action", context.Action())
	}
}
