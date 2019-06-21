package options

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
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
