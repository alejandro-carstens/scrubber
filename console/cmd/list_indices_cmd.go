package cmd

import (
	"scrubber/actions/contexts"
	"scrubber/actions/options"
	"scrubber/logging"

	"github.com/spf13/cobra"
)

type listIndicesCmd struct {
	baseActionCmd
}

func (lic *listIndicesCmd) new(logger *logging.SrvLogger) *cobra.Command {
	command := &cobra.Command{
		Use:   "list-indices",
		Short: "list all indices",
		Args:  lic.Validate,
		Run:   lic.Handle,
	}

	lic.logger = logger

	return command
}

func (lic *listIndicesCmd) Validate(cmd *cobra.Command, args []string) error {
	lic.context = new(contexts.ListIndicesContext)

	options := &options.ListIndicesOptions{}
	options.BindFlags(cmd.Flags())

	return lic.context.SetOptions(options)
}
