package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
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

func (cio *CreateIndexOptions) BindFlags(flags *pflag.FlagSet) error {
	cio.defaultBindFlags(flags)

	cio.Name = stringFromFlags(flags, "name")

	if len(stringFromFlags(flags, "extra_settings")) > 0 {
		cio.ExtraSettings = map[string]interface{}{}

		if err := json.Unmarshal([]byte(stringFromFlags(flags, "extra_settings")), &cio.ExtraSettings); err != nil {
			return err
		}
	}

	return nil
}
