package cmd

import (
	"errors"
	"scrubber/actions/contexts"
	"scrubber/actions/options"
	"scrubber/logger"

	"github.com/spf13/cobra"
)

type closeIndicesCmd struct {
	baseActionCmd
}

func (cic *closeIndicesCmd) new(logger *logger.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:   "close-indices",
		Short: "closes the specified list of indices",
		Args:  cic.Validate,
		Run:   cic.Handle,
	}

	command.Flags().StringSlice("indices", []string{}, "indices to be closed")

	cic.logger = logger

	return command
}

func (cic *closeIndicesCmd) Validate(cmd *cobra.Command, args []string) error {
	indices, _ := cmd.Flags().GetStringSlice("indices")

	if len(indices) == 0 {
		return errors.New("a list of indices is required to run this command")
	}

	cic.context = new(contexts.CloseIndicesContext)
	cic.context.SetList(indices...)

	options := &options.CloseIndicesOptions{}
	options.BindFlags(cmd.Flags())

	return cic.context.SetOptions(options)
}
