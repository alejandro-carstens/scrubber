package resourcepool

import (
	"scrubber/config"

	"github.com/alejandro-carstens/golastic"
)

type elasticsearch struct{}

func (e *elasticsearch) boot(rp *ResourcePool) error {
	configuration, err := config.GetFromEnvVars("elasticsearch")

	if err != nil {
		return err
	}

	cnf := configuration.(*config.Elasticsearch)

	connection := golastic.NewConnection(&golastic.ConnectionContext{
		Urls:                cnf.Urls,
		Password:            cnf.Password,
		Username:            cnf.Username,
		HealthCheckInterval: cnf.HealthCheckInterval,
		ErrorLogPrefix:      cnf.ErrorLogPrefix,
		InfoLogPrefix:       cnf.InfoLogPrefix,
		Context:             rp.Context(),
	})

	if err := connection.Connect(); err != nil {
		return err
	}

	rp.elasticsearch = connection

	return nil
}
