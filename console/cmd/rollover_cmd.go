package cmd

import (
	"scrubber/actions/contexts"
	"scrubber/actions/options"
	rp "scrubber/resourcepool"

	"github.com/spf13/cobra"
)

type rolloverCmd struct {
	baseActionCmd
}

func (rc *rolloverCmd) new() *cobra.Command {
	command := &cobra.Command{
		Use:   "rollover",
		Short: "rollover the specified index alias",
		Args:  rc.Validate,
		Run:   rc.Handle,
	}

	command.Flags().String("name", "", "alias of the index to be rolled over")
	command.Flags().String("max_date", "", "the maximum age of the index")
	command.Flags().Int("max_docs", 0, "the maximum number of documents the index should contain")
	command.Flags().String("max_size", "", "the maximum estimated size of the primary shard of the index")
	command.Flags().String("new_index", "", "the name to be given to the new index")
	command.Flags().String("index_settings", "", "the settings to be provided when creating the new index")

	rc.logger = rp.Logger()

	return command
}

// Validate implementation of the Commandable interface
func (rc *rolloverCmd) Validate(cmd *cobra.Command, args []string) error {
	rc.context = new(contexts.RolloverContext)

	options := &options.RolloverOptions{}

	if err := options.BindFlags(cmd.Flags()); err != nil {
		return err
	}

	return rc.context.SetOptions(options)
}
