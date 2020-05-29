package cmd

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"scrubber/console"
	"scrubber/logger"
	"scrubber/notifications"
	rp "scrubber/resourcepool"

	"github.com/spf13/cobra"
)

const DEFAULT_NOTIFICATIONS_QUEUE_CAPACITY int = 10000

type schedulerCmd struct {
	logger  *logger.Logger
	path    string
	exclude []string
}

func (sc *schedulerCmd) new() *cobra.Command {
	command := &cobra.Command{
		Use:   "scheduler",
		Short: "run the scheduler",
		Args:  sc.Validate,
		Run:   sc.Handle,
	}

	command.Flags().String("path", "", "the path of the directory or file containing the actions config")
	command.Flags().StringSlice("eclude", []string{}, "exclude files that matches any element of the comma separated list of words")

	sc.logger = rp.Logger()

	return command
}

// Validate implementation of the Commandable interface
func (sc *schedulerCmd) Validate(cmd *cobra.Command, args []string) error {
	path, _ := cmd.Flags().GetString("path")

	if len(path) == 0 {
		path = os.Getenv("ACTIONS_PATH")
	}

	if len(path) == 0 {
		return errors.New("the path field is required")
	}

	exclude, _ := cmd.Flags().GetStringSlice("exclude")

	if len(exclude) == 0 {
		exclude = strings.Split(os.Getenv("EXCLUDE_MATCHES"), ",")
	}

	sc.path = path
	sc.exclude = exclude

	return nil
}

func (sc *schedulerCmd) Handle(cmd *cobra.Command, args []string) {
	var capacity int = DEFAULT_NOTIFICATIONS_QUEUE_CAPACITY

	if len(os.Getenv("NOTIFICATIONS_QUEUE_CAPACITY")) > 0 {
		value, err := strconv.Atoi(os.Getenv("NOTIFICATIONS_QUEUE_CAPACITY"))

		if err != nil {
			sc.logger.Debugf(err.Error())

			capacity = int(value)
		}
	}

	queue, err := notifications.NewQueue(capacity, sc.logger)

	if err != nil {
		sc.logger.Errorf(err.Error())

		return
	}

	defer queue.Release()

	if err := console.NewScheduler(sc.path, sc.exclude, queue).Run(); err != nil {
		sc.logger.Errorf(err.Error())
	}
}
