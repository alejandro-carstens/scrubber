package options

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
)

type CloseIndicesOptions struct {
	defaultOptions
}

func (cio *CloseIndicesOptions) FillFromContainer(container *gabs.Container) error {
	cio.container = container

	return json.Unmarshal(container.Bytes(), cio)
}

func (cio *CloseIndicesOptions) Validate() error {
	return nil
}
