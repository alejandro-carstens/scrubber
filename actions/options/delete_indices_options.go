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
	return dio.defaultBindFlags(flags)
}
