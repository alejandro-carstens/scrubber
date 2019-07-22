package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type RolloverOptions struct {
	defaultOptions
	Name          string                 `json:"name"`
	MaxAge        string                 `json:"max_age"`
	MaxDocs       int                    `json:"max_docs"`
	MaxSize       string                 `json:"max_size"`
	ExtraSettings map[string]interface{} `json:"extra_settings"`
	NewIndex      string                 `json:"new_index"`
}

func (ro *RolloverOptions) FillFromContainer(container *gabs.Container) error {
	ro.container = container

	return json.Unmarshal(container.Bytes(), ro)
}

func (ro *RolloverOptions) Validate() error {
	if len(ro.Name) == 0 {
		return errors.New("name is a required field")
	}

	if len(ro.MaxAge) == 0 && ro.MaxDocs == 0 && len(ro.MaxSize) == 0 {
		return errors.New("at least one condition needs to be specified")
	}

	return nil
}

func (ro *RolloverOptions) BindFlags(flags *pflag.FlagSet) error {
	ro.defaultBindFlags(flags)

	ro.Name = stringFromFlags(flags, "name")
	ro.MaxAge = stringFromFlags(flags, "max_date")
	ro.MaxDocs = intFromFlags(flags, "max_docs")
	ro.MaxSize = stringFromFlags(flags, "max_size")
	ro.NewIndex = stringFromFlags(flags, "new_index")

	if len(stringFromFlags(flags, "index_settings")) > 0 {
		ro.ExtraSettings = map[string]interface{}{}

		if err := json.Unmarshal([]byte(stringFromFlags(flags, "index_settings")), &ro.ExtraSettings); err != nil {
			return err
		}
	}

	return nil
}
