package cmd

import (
	"scrubber/logging"

	"github.com/spf13/cobra"
)

func Run(logger *logging.SrvLogger) {
	Init(logger).Execute()
}

func Init(logger *logging.SrvLogger) *cobra.Command {
	rootCmd := &cobra.Command{Use: "scrubber"}

	rootCmd.PersistentFlags().Int("timeout", 300, "elasticsearch operation timeout")
	rootCmd.PersistentFlags().Bool("disable_action", false, "flag for preventing the action to be ran")
	rootCmd.AddCommand(new(createIndexCmd).new(logger))
	rootCmd.AddCommand(new(deleteIndicesCmd).new(logger))
	rootCmd.AddCommand(new(closeIndicesCmd).new(logger))
	rootCmd.AddCommand(new(openIndicesCmd).new(logger))
	rootCmd.AddCommand(new(aliasCmd).new(logger))
	rootCmd.AddCommand(new(createRepositoryCmd).new(logger))
	rootCmd.AddCommand(new(snapshotCmd).new(logger))

	return rootCmd
}