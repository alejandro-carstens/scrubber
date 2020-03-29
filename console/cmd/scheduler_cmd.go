package cmd

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/console"
	"github.com/alejandro-carstens/scrubber/logger"
	"github.com/alejandro-carstens/scrubber/notifications"
	"github.com/spf13/cobra"
)

const DEFAULT_HEALTH_CHECK_INTERVAL int64 = 30
const DEFAULT_NOTIFICATIONS_QUEUE_CAPACITY int = 10000

type schedulerCmd struct {
	logger  *logger.Logger
	path    string
	exclude []string
}

func (sc *schedulerCmd) new(logger *logger.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:   "scheduler",
		Short: "run the scheduler",
		Args:  sc.Validate,
		Run:   sc.Handle,
	}

	command.Flags().String("path", "", "the path of the directory or file containing the actions config")
	command.Flags().StringSlice("eclude", []string{}, "exclude files that matches any element of the comma separated list of words")

	sc.logger = logger

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
	var healthCheckInterval int64 = DEFAULT_HEALTH_CHECK_INTERVAL

	if len(os.Getenv("ELASTICSEARCH_HEALTH_CHECK_INTERVAL")) > 0 {
		value, err := strconv.Atoi(os.Getenv("ELASTICSEARCH_HEALTH_CHECK_INTERVAL"))

		if err != nil {
			sc.logger.Errorf(err.Error())

			return
		}

		healthCheckInterval = int64(value)
	}

	ctx := context.Background()

	connection := golastic.NewConnection(&golastic.ConnectionContext{
		Urls:                []string{os.Getenv("ELASTICSEARCH_URI")},
		Password:            os.Getenv("ELASTICSEARCH_PASSWORD"),
		Username:            os.Getenv("ELASTICSEARCH_USERNAME"),
		HealthCheckInterval: healthCheckInterval,
		ErrorLogPrefix:      os.Getenv("ELASTICSEARCH_ERROR_LOG_PREFIX"),
		InfoLogPrefix:       os.Getenv("ELASTICSEARCH_INFO_LOG_PREFIX"),
		Context:             ctx,
	})

	if err := connection.Connect(); err != nil {
		sc.logger.Errorf("%v [ELASTICSEARCH_URI: %v]", err.Error(), os.Getenv("ELASTICSEARCH_URI"))

		return
	}

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

	if err := console.NewScheduler(sc.path, sc.exclude, sc.logger, connection, queue, ctx).Run(); err != nil {
		sc.logger.Errorf(err.Error())
	}
}
