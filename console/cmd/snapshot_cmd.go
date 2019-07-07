package cmd

import (
	"errors"
	"scrubber/actions/contexts"
	"scrubber/actions/options"
	"scrubber/logger"

	"github.com/spf13/cobra"
)

type snapshotCmd struct {
	baseActionCmd
}

func (sc *snapshotCmd) new(logger *logger.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:   "snapshot",
		Short: "take a snapshot of the specified list of indices",
		Args:  sc.Validate,
		Run:   sc.Handle,
	}

	command.Flags().StringSlice("indices", []string{}, "indices to be snapshotted")
	command.Flags().String("repository", "", "the snapshot repository, this field is required")
	command.Flags().String("name", "", "the name of the snapshot")
	command.Flags().Bool("ignore_unavailable", false, "if false and an index is missing the snapshot request will fail")
	command.Flags().Bool("include_global_state", false, "whether Elasticsearch should include the global cluster state with the snapshot")
	command.Flags().Bool("partial", false, "if true, the snapshot will fail if one or more indices being added to the snapshot do not have all primary shards available")
	command.Flags().Bool("wait_for_completion", false, "whether or not the request should return immediately or wait for the operation to complete before returning")
	command.Flags().Int("max_wait", 0, "specifies how long in seconds to wait to see if the action has completed before giving up")
	command.Flags().Int("wait_interval", 0, "specifies how long in seconds to wait to see if the action has completed before giving up")

	sc.logger = logger

	return command
}

func (sc *snapshotCmd) Validate(cmd *cobra.Command, args []string) error {
	indices, _ := cmd.Flags().GetStringSlice("indices")

	if len(indices) == 0 {
		return errors.New("a list of indices is required to run this command")
	}

	sc.context = new(contexts.SnapshotContext)
	sc.context.SetList(indices...)

	options := &options.SnapshotOptions{}
	options.BindFlags(cmd.Flags())

	return sc.context.SetOptions(options)
}
