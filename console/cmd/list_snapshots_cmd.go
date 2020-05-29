package cmd

import (
	"scrubber/actions/contexts"
	"scrubber/actions/options"
	rp "scrubber/resourcepool"

	"github.com/spf13/cobra"
)

type listSnapshotsCmd struct {
	baseActionCmd
}

func (lsc *listSnapshotsCmd) new() *cobra.Command {
	command := &cobra.Command{
		Use:   "list-snapshots",
		Short: "list all snapshots for the given repository",
		Args:  lsc.Validate,
		Run:   lsc.Handle,
	}

	command.Flags().String("repository", "", "the snapshots repository, this field is required")

	lsc.logger = rp.Logger()

	return command
}

// Validate implementation of the Commandable interface
func (lsc *listSnapshotsCmd) Validate(cmd *cobra.Command, args []string) error {
	lsc.context = new(contexts.ListSnapshotsContext)

	options := &options.ListSnapshotsOptions{}
	options.BindFlags(cmd.Flags())

	return lsc.context.SetOptions(options)
}
