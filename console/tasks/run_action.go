package tasks

import (
	"log"
	"scrubber/actions"
	"scrubber/actions/contexts"
	"scrubber/logging"

	"github.com/Jeffail/gabs"
)

type RunAction struct {
	config *gabs.Container
}

func (rat *RunAction) SetConfig(config *gabs.Container) *RunAction {
	rat.config = config

	return rat
}

func (rat *RunAction) Execute() {
	context, err := contexts.New(rat.config)

	if err != nil {
		log.Println(err)

		return
	}

	action, err := actions.Create(context, logging.NewSrvLogger("", true, true, true, true))

	if err != nil {
		log.Println(err)

		return
	}

	if action.DisableAction() {
		log.Println("Action disabled")

		return
	}

	if action.Perform().HasErrors() {
		log.Println("Something went wrong while performing the action. Please check the logs and try again")

		return
	}
}
