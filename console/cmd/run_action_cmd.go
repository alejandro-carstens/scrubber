package cmd

import (
	"errors"
	"os"

	"github.com/alejandro-carstens/scrubber/actions/contexts"
	"github.com/alejandro-carstens/scrubber/logger"
	"github.com/alejandro-carstens/scrubber/ymlparser"
	"github.com/spf13/cobra"
)

type runActionCmd struct {
	baseActionCmd
}

func (rac *runActionCmd) new(logger *logger.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:   "run-action",
		Short: "runs an action from the specified action file",
		Args:  rac.Validate,
		Run:   rac.Handle,
	}

	command.Flags().String("file_path", "", "file path relative to scrubber's root path")

	rac.logger = logger

	return command
}

// Validate implementation of the Commandable interface
func (rac *runActionCmd) Validate(cmd *cobra.Command, args []string) error {
	filePath, err := cmd.Flags().GetString("file_path")

	if err != nil {
		return err
	}

	if len(filePath) == 0 {
		return errors.New("file path is a required field")
	}

	currentDirectory, err := os.Getwd()

	if err != nil {
		return err
	}

	config, err := ymlparser.Parse(currentDirectory + filePath)

	if err != nil {
		return err
	}

	context, err := contexts.New(config)

	if err != nil {
		return err
	}

	rac.context = context

	return nil
}
