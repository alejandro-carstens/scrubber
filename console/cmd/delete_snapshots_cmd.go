package cmd

import (
	"errors"
	"github.com/alejandro-carstens/scrubber/actions/contexts"
	"github.com/alejandro-carstens/scrubber/actions/options"
	"github.com/alejandro-carstens/scrubber/logger"

	"github.com/spf13/cobra"
)

type deleteSnaphotsCmd struct {
	baseActionCmd
}

func (dsc *deleteSnaphotsCmd) new(logger *logger.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:   "delete-snapshots",
		Short: "delete the specified snapshots",
		Args:  dsc.Validate,
		Run:   dsc.Handle,
	}

	command.Flags().StringSlice("snapshots", []string{}, "snapshots to be deleted")
	command.Flags().String("repository", "", "the snapshot repository, this field is required")
	command.Flags().Int("retry_count", 0, "number of times the delete call should happen if it fails")
	command.Flags().Int("retry_interval", 0, "number of seconds the delete call should wait before it is retried")

	dsc.logger = logger

	return command
}

func (dsc *deleteSnaphotsCmd) Validate(cmd *cobra.Command, args []string) error {
	snapshots, _ := cmd.Flags().GetStringSlice("snapshots")

	if len(snapshots) == 0 {
		return errors.New("a list of snapshots is required to run this command")
	}

	dsc.context = new(contexts.DeleteSnapshotsContext)
	dsc.context.SetList(snapshots...)

	options := &options.DeleteSnapshotsOptions{}
	options.BindFlags(cmd.Flags())

	return dsc.context.SetOptions(options)
}
