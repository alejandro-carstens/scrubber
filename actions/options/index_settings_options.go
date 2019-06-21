package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
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
