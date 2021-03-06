package cmd

import (
	"os"

	"github.com/alejandro-carstens/scrubber/logger"
	"github.com/spf13/cobra"
)

func Run() error {
	return boot(logger.NewLogger(os.Getenv("LOG_FILE"), true, true, true, true)).Execute()
}

func boot(logger *logger.Logger) *cobra.Command {
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
	scrubber.AddCommand(new(listIndicesCmd).new(logger))
	scrubber.AddCommand(new(listSnapshotsCmd).new(logger))
	scrubber.AddCommand(new(schedulerCmd).new(logger))
	scrubber.AddCommand(new(deleteRepositoriesCmd).new(logger))
	scrubber.AddCommand(new(rolloverCmd).new(logger))

	return scrubber
}
