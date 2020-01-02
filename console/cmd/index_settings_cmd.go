package cmd

import (
	"errors"

	"github.com/alejandro-carstens/scrubber/actions/contexts"
	"github.com/alejandro-carstens/scrubber/actions/options"
	"github.com/alejandro-carstens/scrubber/logger"
	"github.com/spf13/cobra"
)

type indexSettingsCmd struct {
	baseActionCmd
}

func (isc *indexSettingsCmd) new(logger *logger.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:   "index-settings",
		Short: "updates the index settings for the indices in the actionable list",
		Args:  isc.Validate,
		Run:   isc.Handle,
	}

	isc.logger = logger

	command.Flags().StringSlice("indices", []string{}, "list of indices whose settings need to be updated")
	command.Flags().String("index_settings", "", "index settings to be applied as JSON")

	return command
}

// Validate implementation of the Commandable interface
func (isc *indexSettingsCmd) Validate(cmd *cobra.Command, args []string) error {
	indices, _ := cmd.Flags().GetStringSlice("indices")

	if len(indices) == 0 {
		return errors.New("a list of indices is required to run this command")
	}

	isc.context = new(contexts.IndexSettingsContext)
	isc.context.SetList(indices...)

	options := &options.IndexSettingsOptions{}

	if err := options.BindFlags(cmd.Flags()); err != nil {
		return err
	}

	return isc.context.SetOptions(options)
}
