package cmd

import (
	"errors"
	"os"
	"scrubber/console"
	"scrubber/logger"
	"strconv"

	"github.com/alejandro-carstens/golastic"
	"github.com/spf13/cobra"
)

const DEFAULT_HEALTH_CHECK_INTERVAL int64 = 30

type schedulerCmd struct {
	logger *logger.Logger
	path   string
}

func (sc *schedulerCmd) new(logger *logger.Logger) *cobra.Command {
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
		path = os.Getenv("ACTIONS_PATH")
	}

	if len(path) == 0 {
		return errors.New("the path field is required")
	}

	sc.path = path

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

	connection := golastic.NewConnection(&golastic.ConnectionContext{
		Urls:                []string{os.Getenv("ELASTICSEARCH_URI")},
		Password:            os.Getenv("ELASTICSEARCH_PASSWORD"),
		Username:            os.Getenv("ELASTICSEARCH_USERNAME"),
		HealthCheckInterval: healthCheckInterval,
		ErrorLogPrefix:      os.Getenv("ELASTICSEARCH_ERROR_LOG_PREFIX"),
		InfoLogPrefix:       os.Getenv("ELASTICSEARCH_INFO_LOG_PREFIX"),
	})

	if err := connection.Connect(); err != nil {
		sc.logger.Errorf(err.Error())

		return
	}

	if err := console.NewScheduler(sc.path, sc.logger, connection).Run(); err != nil {
		sc.logger.Errorf(err.Error())
	}
}
