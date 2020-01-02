package cmd

import (
	"github.com/alejandro-carstens/scrubber/actions/contexts"
	"github.com/alejandro-carstens/scrubber/console"
	"github.com/alejandro-carstens/scrubber/logger"
	"github.com/spf13/cobra"
)

type baseActionCmd struct {
	context contexts.Contextable
	logger  *logger.Logger
}

// Handle implementation of the Commandable interface
func (bac *baseActionCmd) Handle(cmd *cobra.Command, args []string) {
	console.Execute(bac.context, bac.logger, nil, nil)
}
