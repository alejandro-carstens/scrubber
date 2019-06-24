package options

import (
	"encoding/json"
	"errors"

	"github.com/spf13/pflag"

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

func (cio *CreateIndexOptions) BindFlags(flags *pflag.FlagSet) error {
	cio.defaultBindFlags(flags)

	name, _ := flags.GetString("name")
	extraSettings, _ := flags.GetString("extra_settings")

	cio.Name = name

	if len(extraSettings) > 0 {
		es := map[string]interface{}{}

		if err := json.Unmarshal([]byte(extraSettings), &es); err != nil {
			return err
		}

		cio.ExtraSettings = es
	}

	return nil
}
