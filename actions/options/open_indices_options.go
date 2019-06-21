package options

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
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
