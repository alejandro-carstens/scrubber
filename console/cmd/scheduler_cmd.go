package cmd

import (
	"errors"
	"scrubber/console"
	"scrubber/logging"

	"github.com/spf13/cobra"
)

type schedulerCmd struct {
	logger *logging.SrvLogger
	path   string
}

func (sc *schedulerCmd) new(logger *logging.SrvLogger) *cobra.Command {
	command := &cobra.Command{
		Use:   "scheduler",
		Short: "run the scheduler",
		Args:  sc.Validate,
		Run:   sc.Handle,
	}

	command.Flags().String("path", "", "the path of the directory or file containing the actions config")

	sc.logger = logger

	return command
}

func (sc *schedulerCmd) Validate(cmd *cobra.Command, args []string) error {
	path, _ := cmd.Flags().GetString("path")

	if len(path) == 0 {
		return errors.New("the path field is required")
	}

	sc.path = path

	return nil
}

func (sc *schedulerCmd) Handle(cmd *cobra.Command, args []string) {
	if err := console.NewScheduler(sc.path, sc.logger).Run(); err != nil {
		sc.logger.Errorf(err.Error())
	}
}
