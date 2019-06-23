package commands

import (
	"scrubber/actions/contexts"
	"scrubber/console"
	"scrubber/logging"

	"github.com/spf13/cobra"
)

type baseCommand struct {
	context contexts.Contextable
	logger  *logging.SrvLogger
}

func (bc *baseCommand) Handle(cmd *cobra.Command, args []string) {
	console.Execute(bc.context, bc.logger)
}
