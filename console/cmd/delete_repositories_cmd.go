package cmd

import (
	"scrubber/actions/contexts"
	"scrubber/actions/options"
	rp "scrubber/resourcepool"

	"github.com/spf13/cobra"
)

type deleteRepositoriesCmd struct {
	baseActionCmd
}

func (drc *deleteRepositoriesCmd) new() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete-repositories",
		Short: "delete the specified list of repositories",
		Args:  drc.Validate,
		Run:   drc.Handle,
	}

	command.Flags().String("repositories", "", "repositories to be deleted")

	drc.logger = rp.Logger()

	return command
}

// Validate implementation of the Commandable interface
func (drc *deleteRepositoriesCmd) Validate(cmd *cobra.Command, args []string) error {
	drc.context = new(contexts.DeleteRepositoriesContext)

	options := &options.DeleteRepositoriesOptions{}
	options.BindFlags(cmd.Flags())

	return drc.context.SetOptions(options)
}
