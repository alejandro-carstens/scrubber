package cmd

import (
	"errors"

	"scrubber/actions/contexts"
	"scrubber/actions/options"
	rp "scrubber/resourcepool"

	"github.com/spf13/cobra"
)

type deleteIndicesCmd struct {
	baseActionCmd
}

func (dic *deleteIndicesCmd) new() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete-indices",
		Short: "deletes the specified list of indices",
		Args:  dic.Validate,
		Run:   dic.Handle,
	}

	command.Flags().StringSlice("indices", []string{}, "indices to be deleted")

	dic.logger = rp.Logger()

	return command
}

// Validate implementation of the Commandable interface
func (dic *deleteIndicesCmd) Validate(cmd *cobra.Command, args []string) error {
	indices, _ := cmd.Flags().GetStringSlice("indices")

	if len(indices) == 0 {
		return errors.New("a list of indices is required to run this command")
	}

	dic.context = new(contexts.DeleteIndicesContext)
	dic.context.SetList(indices...)

	options := &options.DeleteIndicesOptions{}
	options.BindFlags(cmd.Flags())

	return dic.context.SetOptions(options)
}
