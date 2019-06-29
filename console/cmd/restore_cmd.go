package cmd

import (
	"scrubber/actions/contexts"
	"scrubber/actions/options"
	"scrubber/logging"

	"github.com/spf13/cobra"
)

type restoreCmd struct {
	baseActionCmd
}

func (rc *restoreCmd) new(logger *logging.SrvLogger) *cobra.Command {
	command := &cobra.Command{
		Use:   "restore",
		Short: "restore indices for a specifc snapshot",
		Args:  rc.Validate,
		Run:   rc.Handle,
	}

	command.Flags().String("name", "", "name of the snapshot to restore")
	command.Flags().String("indices", "", "comma separated list of indices to restore")
	command.Flags().String("repository", "", "the repository of the snapshots being restored, this field is required")
	command.Flags().String("rename_pattern", "", "used to rename indices on restore using regular expression that supports referencing the original index name")
	command.Flags().String("rename_replacement", "", "used to rename indices on restore using regular expression that supports referencing the original index name")
	command.Flags().String("extra_settings", "", "index settings for the indices being restored in JSON format")
	command.Flags().Int("max_wait", 0, "specifies how long in seconds to wait to see if the action has completed before giving up")
	command.Flags().Int("wait_interval", 0, "specifies how long in seconds to wait to see if the action has completed before giving up")
	command.Flags().Bool("ignore_unavailable", false, "if false and an index is missing the snapshot request will fail")
	command.Flags().Bool("include_global_state", false, "whether Elasticsearch should include the global cluster state with the snapshot")
	command.Flags().Bool("partial", false, "if true, the snapshot will fail if one or more indices being added to the snapshot do not have all primary shards available")
	command.Flags().Bool("wait_for_completion", false, "whether or not the request should return immediately or wait for the operation to complete before returning")
	command.Flags().Bool("include_aliases", false, "")

	rc.logger = logger

	return command
}

func (rc *restoreCmd) Validate(cmd *cobra.Command, args []string) error {
	rc.context = new(contexts.RestoreContext)

	options := &options.DeleteSnapshotsOptions{}

	if err := options.BindFlags(cmd.Flags()); err != nil {
		return err
	}

	return rc.context.SetOptions(options)
}
