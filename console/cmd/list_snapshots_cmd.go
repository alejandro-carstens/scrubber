package cmd

import (
	"scrubber/actions/contexts"
	"scrubber/actions/options"
	"scrubber/logging"

	"github.com/spf13/cobra"
)

type listSnapshotsCmd struct {
	baseActionCmd
}

func (lsc *listSnapshotsCmd) new(logger *logging.SrvLogger) *cobra.Command {
	command := &cobra.Command{
		Use:   "list-snapshots",
		Short: "list all snapshots for the given repository",
		Args:  lsc.Validate,
		Run:   lsc.Handle,
	}

	command.Flags().String("repository", "", "the snapshots repository, this field is required")

	lsc.logger = logger

	return command
}

func (lsc *listSnapshotsCmd) Validate(cmd *cobra.Command, args []string) error {
	lsc.context = new(contexts.ListSnapshotsContext)

	options := &options.ListSnapshotsOptions{}
	options.BindFlags(cmd.Flags())

	return lsc.context.SetOptions(options)
}
