package cmd

import (
	"github.com/alejandro-carstens/scrubber/actions/contexts"
	"github.com/alejandro-carstens/scrubber/actions/options"
	"github.com/alejandro-carstens/scrubber/logger"

	"github.com/spf13/cobra"
)

type deleteRepositoriesCmd struct {
	baseActionCmd
}

func (drc *deleteRepositoriesCmd) new(logger *logger.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:   "delete-repositories",
		Short: "delete the specified list of repositories",
		Args:  drc.Validate,
		Run:   drc.Handle,
	}

	command.Flags().String("repositories", "", "repositories to be deleted")

	drc.logger = logger

	return command
}

func (drc *deleteRepositoriesCmd) Validate(cmd *cobra.Command, args []string) error {
	drc.context = new(contexts.DeleteRepositoriesContext)

	options := &options.DeleteRepositoriesOptions{}
	options.BindFlags(cmd.Flags())

	return drc.context.SetOptions(options)
}
