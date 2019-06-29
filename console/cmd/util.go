package cmd

import (
	"scrubber/logging"

	"github.com/spf13/cobra"
)

func Run(logger *logging.SrvLogger) {
	Init(logger).Execute()
}

func Init(logger *logging.SrvLogger) *cobra.Command {
	actionCmd := &cobra.Command{Use: "scrubber"}

	actionCmd.PersistentFlags().Int("timeout", 300, "elasticsearch operation timeout")
	actionCmd.PersistentFlags().Bool("disable_action", false, "flag for preventing the action to be ran")
	actionCmd.AddCommand(new(createIndexCmd).new(logger))
	actionCmd.AddCommand(new(deleteIndicesCmd).new(logger))
	actionCmd.AddCommand(new(closeIndicesCmd).new(logger))
	actionCmd.AddCommand(new(openIndicesCmd).new(logger))
	actionCmd.AddCommand(new(aliasCmd).new(logger))
	actionCmd.AddCommand(new(createRepositoryCmd).new(logger))
	actionCmd.AddCommand(new(snapshotCmd).new(logger))
	actionCmd.AddCommand(new(deleteSnaphotsCmd).new(logger))
	actionCmd.AddCommand(new(restoreCmd).new(logger))
	actionCmd.AddCommand(new(indexSettingsCmd).new(logger))

	return actionCmd
}
