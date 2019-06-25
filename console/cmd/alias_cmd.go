package cmd

import (
	"errors"
	"scrubber/actions/contexts"
	"scrubber/actions/options"
	"scrubber/logging"

	"github.com/spf13/cobra"
)

type aliasCmd struct {
	baseActionCmd
}

func (ac *aliasCmd) new(logger *logging.SrvLogger) *cobra.Command {
	command := &cobra.Command{
		Use:   "alias",
		Short: "add or remove indices from an alias",
		Args:  ac.Validate,
		Run:   ac.Handle,
	}

	command.Flags().StringSlice("indices", []string{}, "indices to be closed")
	command.Flags().String("name", "", "alias name")
	command.Flags().String("type", "", "alias action type [add or remove]")
	command.Flags().String("filter", "", "alias filter extra setting in JSON format")
	command.Flags().String("routing", "", "alias routing extra setting")
	command.Flags().String("search_routing", "", "alias search routing extra setting")

	ac.logger = logger

	return command
}

func (ac *aliasCmd) Validate(cmd *cobra.Command, args []string) error {
	indices, _ := cmd.Flags().GetStringSlice("indices")

	if len(indices) == 0 {
		return errors.New("a list of indices is required to run this command")
	}

	ac.context = new(contexts.AliasContext)
	ac.context.SetList(indices...)

	options := &options.AliasOptions{}

	if err := options.BindFlags(cmd.Flags()); err != nil {
		return err
	}

	return ac.context.SetOptions(options)
}
