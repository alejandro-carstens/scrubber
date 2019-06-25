package options

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type OpenIndicesOptions struct {
	defaultOptions
}

func (oio *OpenIndicesOptions) FillFromContainer(container *gabs.Container) error {
	oio.container = container

	return json.Unmarshal(container.Bytes(), oio)
}

func (oio *OpenIndicesOptions) Validate() error {
	return nil
}

func (oio *OpenIndicesOptions) BindFlags(flags *pflag.FlagSet) error {
	return oio.defaultBindFlags(flags)
}
