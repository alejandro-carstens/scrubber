package cmd

import (
	"bytes"
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

	return rootCmd
}

func executeCommand(root *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)

	root.SetOutput(buf)

	root.SetArgs(args)

	_, err := root.ExecuteC()

	return buf.String(), err
}
