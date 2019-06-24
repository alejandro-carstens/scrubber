package cmd

import (
	"errors"
	"scrubber/actions/contexts"
	"scrubber/actions/options"
	"scrubber/logging"

	"github.com/spf13/cobra"
)

type deleteIndicesCmd struct {
	baseActionCmd
}

func (dic *deleteIndicesCmd) new(logger *logging.SrvLogger) *cobra.Command {
	command := &cobra.Command{
		Use:   "delete-indices",
		Short: "deletes the specified list of indices",
		Args:  dic.Validate,
		Run:   dic.Handle,
	}

	command.Flags().Int("timeout", 300, "elasticsearch operation timeout")
	command.Flags().Bool("disable_action", false, "flag for preventing the action to be ran")
	command.Flags().StringSlice("indices", []string{}, "indices to be deleted")

	dic.logger = logger

	return command
}

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
