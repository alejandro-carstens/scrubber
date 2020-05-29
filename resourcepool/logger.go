package resourcepool

import (
	"scrubber/config"
	logs "scrubber/logger"
)

type logger struct{}

func (l *logger) boot(rp *ResourcePool) error {
	configuration, err := config.GetFromEnvVars("logger")

	if err != nil {
		return err
	}

	cnf := configuration.(*config.Logger)

	rp.logger = logs.NewLogger(cnf.LogFile, true, true, true, true)

	return nil
}
