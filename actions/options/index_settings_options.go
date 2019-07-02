package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type IndexSettingsOptions struct {
	defaultOptions
	IndexSettings map[string]interface{} `json:"index_settings"`
}

func (iso *IndexSettingsOptions) FillFromContainer(container *gabs.Container) error {
	iso.container = container

	return json.Unmarshal(container.Bytes(), iso)
}

func (iso *IndexSettingsOptions) Validate() error {
	if iso.IndexSettings == nil {
		return errors.New("index_settings option is required")
	}

	return nil
}

func (iso *IndexSettingsOptions) BindFlags(flags *pflag.FlagSet) error {
	iso.defaultBindFlags(flags)

	if len(stringFromFlags(flags, "index_settings")) > 0 {
		iso.IndexSettings = map[string]interface{}{}

		if err := json.Unmarshal([]byte(stringFromFlags(flags, "index_settings")), &iso.IndexSettings); err != nil {
			return err
		}
	}

	return nil
}
