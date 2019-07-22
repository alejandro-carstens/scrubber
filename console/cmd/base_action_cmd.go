package cmd

import (
	"scrubber/actions/contexts"
	"scrubber/console"
	"scrubber/logger"

	"github.com/spf13/cobra"
)

type baseActionCmd struct {
	context contexts.Contextable
	logger  *logger.Logger
}

func (bac *baseActionCmd) Handle(cmd *cobra.Command, args []string) {
	console.Execute(bac.context, bac.logger, nil)
}
