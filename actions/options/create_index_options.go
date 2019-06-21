package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
)

type CreateIndexOptions struct {
	defaultOptions
	Name          string                 `json:"name"`
	ExtraSettings map[string]interface{} `json:"extra_settings"`
}

func (cio *CreateIndexOptions) FillFromContainer(container *gabs.Container) error {
	cio.container = container

	return json.Unmarshal(container.Bytes(), cio)
}

func (cio *CreateIndexOptions) Validate() error {
	if len(cio.Name) == 0 {
		return errors.New("The name option is required")
	}

	return nil
}
