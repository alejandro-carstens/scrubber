package cmd

import (
	"scrubber/logging"

	"github.com/spf13/cobra"
)

func Run(logger *logging.SrvLogger) {
	Init(logger).Execute()
}

func Init(logger *logging.SrvLogger) *cobra.Command {
	scrubber := &cobra.Command{Use: "scrubber"}

	scrubber.PersistentFlags().Int("timeout", 300, "elasticsearch operation timeout")
	scrubber.PersistentFlags().Bool("disable_action", false, "flag for preventing the action to be ran")

	scrubber.AddCommand(new(createIndexCmd).new(logger))
	scrubber.AddCommand(new(deleteIndicesCmd).new(logger))
	scrubber.AddCommand(new(closeIndicesCmd).new(logger))
	scrubber.AddCommand(new(openIndicesCmd).new(logger))
	scrubber.AddCommand(new(aliasCmd).new(logger))
	scrubber.AddCommand(new(createRepositoryCmd).new(logger))
	scrubber.AddCommand(new(snapshotCmd).new(logger))
	scrubber.AddCommand(new(deleteSnaphotsCmd).new(logger))
	scrubber.AddCommand(new(restoreCmd).new(logger))
	scrubber.AddCommand(new(indexSettingsCmd).new(logger))
	scrubber.AddCommand(new(runActionCmd).new(logger))

	return scrubber
}
