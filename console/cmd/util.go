package cmd

import (
	rp "scrubber/resourcepool"

	"github.com/spf13/cobra"
)

func Run() error {
	return boot().Execute()
}

func boot() *cobra.Command {
	rp.Boot("mysql")

	scrubber := &cobra.Command{Use: "scrubber"}

	scrubber.PersistentFlags().Int("timeout", 300, "elasticsearch operation timeout")
	scrubber.PersistentFlags().Bool("disable_action", false, "flag for preventing the action to be ran")

	scrubber.AddCommand(new(createIndexCmd).new())
	scrubber.AddCommand(new(deleteIndicesCmd).new())
	scrubber.AddCommand(new(closeIndicesCmd).new())
	scrubber.AddCommand(new(openIndicesCmd).new())
	scrubber.AddCommand(new(aliasCmd).new())
	scrubber.AddCommand(new(createRepositoryCmd).new())
	scrubber.AddCommand(new(snapshotCmd).new())
	scrubber.AddCommand(new(deleteSnaphotsCmd).new())
	scrubber.AddCommand(new(restoreCmd).new())
	scrubber.AddCommand(new(indexSettingsCmd).new())
	scrubber.AddCommand(new(runActionCmd).new())
	scrubber.AddCommand(new(listIndicesCmd).new())
	scrubber.AddCommand(new(listSnapshotsCmd).new())
	scrubber.AddCommand(new(schedulerCmd).new())
	scrubber.AddCommand(new(deleteRepositoriesCmd).new())
	scrubber.AddCommand(new(rolloverCmd).new())

	return scrubber
}
