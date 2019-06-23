package commands

import (
	"scrubber/actions/contexts"
	"scrubber/console"
	"scrubber/logging"

	"github.com/spf13/cobra"
)

type baseActionCommand struct {
	context contexts.Contextable
	logger  *logging.SrvLogger
}

func (bac *baseActionCommand) Handle(cmd *cobra.Command, args []string) {
	console.Execute(bac.context, bac.logger)
}
