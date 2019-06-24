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

	rootCmd.AddCommand(new(deleteIndicesCmd).new(logger))
	rootCmd.AddCommand(new(createIndexCmd).new(logger))

	return rootCmd
}

func executeCommand(root *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)

	root.SetOutput(buf)

	root.SetArgs(args)

	_, err := root.ExecuteC()

	return buf.String(), err
}
