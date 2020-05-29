package cmd

import (
	"scrubber/actions/contexts"
	"scrubber/actions/options"
	rp "scrubber/resourcepool"

	"github.com/spf13/cobra"
)

type listIndicesCmd struct {
	baseActionCmd
}

func (lic *listIndicesCmd) new() *cobra.Command {
	command := &cobra.Command{
		Use:   "list-indices",
		Short: "list all indices",
		Args:  lic.Validate,
		Run:   lic.Handle,
	}

	lic.logger = rp.Logger()

	return command
}

// Validate implementation of the Commandable interface
func (lic *listIndicesCmd) Validate(cmd *cobra.Command, args []string) error {
	lic.context = new(contexts.ListIndicesContext)

	options := &options.ListIndicesOptions{}
	options.BindFlags(cmd.Flags())

	return lic.context.SetOptions(options)
}
