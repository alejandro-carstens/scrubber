package commands

import (
	"scrubber/logging"

	"github.com/spf13/cobra"
)

func Run() {
	rootCommand := &cobra.Command{Use: "scrubber"}

	logger := logging.NewSrvLogger("", true, true, true, true)

	rootCommand.AddCommand(new(deleteIndicesCommand).new(logger))

	rootCommand.Execute()
}
