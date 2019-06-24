package cmd

import (
	"scrubber/actions/contexts"
	"scrubber/actions/options"
	"scrubber/logging"

	"github.com/spf13/cobra"
)

type createIndexCmd struct {
	baseActionCmd
}

func (cic *createIndexCmd) new(logger *logging.SrvLogger) *cobra.Command {
	command := &cobra.Command{
		Use:   "create-index",
		Short: "creates an index",
		Args:  cic.Validate,
		Run:   cic.Handle,
	}

	command.Flags().Int("timeout", 300, "elasticsearch operation timeout")
	command.Flags().Bool("disable_action", false, "flag for preventing the action to be ran")
	command.Flags().String("name", "", "name of the index to be created")
	command.Flags().String("extra_settings", "", "index settings and mappings as JSON")

	cic.logger = logger

	return command
}

func (cic *createIndexCmd) Validate(cmd *cobra.Command, args []string) error {
	cic.context = new(contexts.CreateIndexContext)

	options := &options.CreateIndexOptions{}

	if err := options.BindFlags(cmd.Flags()); err != nil {
		return err
	}

	return cic.context.SetOptions(options)
}
