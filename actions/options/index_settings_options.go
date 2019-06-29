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

	indexSettings, _ := flags.GetString("index_settings")

	if len(indexSettings) == 0 {
		return nil
	}

	is := map[string]interface{}{}

	if err := json.Unmarshal([]byte(indexSettings), &is); err != nil {
		return err
	}

	iso.IndexSettings = is

	return nil
}
