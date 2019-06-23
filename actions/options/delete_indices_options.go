package options

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type DeleteIndicesOptions struct {
	defaultOptions
}

func (dio *DeleteIndicesOptions) FillFromContainer(container *gabs.Container) error {
	dio.container = container

	return json.Unmarshal(container.Bytes(), dio)
}

func (dio *DeleteIndicesOptions) Validate() error {
	return nil
}

func (dio *DeleteIndicesOptions) BindFlags(flags *pflag.FlagSet) error {
	timeout, _ := flags.GetInt("timeout")
	disableAction, _ := flags.GetBool("disable_action")

	dio.Timeout = timeout
	dio.DisableAction = disableAction

	return nil
}
