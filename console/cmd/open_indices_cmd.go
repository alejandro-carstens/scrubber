package cmd

import (
	"errors"

	"scrubber/actions/contexts"
	"scrubber/actions/options"
	rp "scrubber/resourcepool"

	"github.com/spf13/cobra"
)

type openIndicesCmd struct {
	baseActionCmd
}

func (oic *openIndicesCmd) new() *cobra.Command {
	command := &cobra.Command{
		Use:   "open-indices",
		Short: "open the specified list of indices",
		Args:  oic.Validate,
		Run:   oic.Handle,
	}

	command.Flags().StringSlice("indices", []string{}, "indices to be openned")

	oic.logger = rp.Logger()

	return command
}

// Validate implementation of the Commandable interface
func (oic *openIndicesCmd) Validate(cmd *cobra.Command, args []string) error {
	indices, _ := cmd.Flags().GetStringSlice("indices")

	if len(indices) == 0 {
		return errors.New("a list of indices is required to run this command")
	}

	oic.context = new(contexts.OpenIndicesContext)
	oic.context.SetList(indices...)

	options := &options.OpenIndicesOptions{}
	options.BindFlags(cmd.Flags())

	return oic.context.SetOptions(options)
}
