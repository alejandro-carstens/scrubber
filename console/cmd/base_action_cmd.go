package cmd

import (
	"scrubber/actions/contexts"
	"scrubber/console"
	"scrubber/logging"

	"github.com/spf13/cobra"
)

type baseActionCmd struct {
	context contexts.Contextable
	logger  *logging.SrvLogger
}

func (bac *baseActionCmd) Handle(cmd *cobra.Command, args []string) {
	console.Execute(bac.context, bac.logger)
}
