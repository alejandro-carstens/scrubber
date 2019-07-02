package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type AliasOptions struct {
	defaultOptions
	ExtraSettings *AliasExtraSettingsOption `json:"extra_settings"`
	Name          string                    `json:"name"`
	Type          string                    `json:"type"`
}

type AliasExtraSettingsOption struct {
	Routing       string                 `json:"routing"`
	SearchRouting string                 `json:"search_routing"`
	Filter        map[string]interface{} `json:"filter"`
}

func (ao *AliasOptions) FillFromContainer(container *gabs.Container) error {
	ao.container = container

	return json.Unmarshal(container.Bytes(), ao)
}

func (ao *AliasOptions) Validate() error {
	if len(ao.Name) == 0 {
		return errors.New("name option is required")
	}

	if ao.Type != "add" && ao.Type != "remove" {
		return errors.New("type option can only be add or remove")
	}

	return nil
}

func (ao *AliasOptions) BindFlags(flags *pflag.FlagSet) error {
	ao.defaultBindFlags(flags)

	aliasfilter := map[string]interface{}{}

	if len(stringFromFlags(flags, "filter")) > 0 {
		if err := json.Unmarshal([]byte(stringFromFlags(flags, "filter")), &aliasfilter); err != nil {
			return err
		}
	}

	ao.Name = stringFromFlags(flags, "name")
	ao.Type = stringFromFlags(flags, "type")
	ao.ExtraSettings = &AliasExtraSettingsOption{
		Routing:       stringFromFlags(flags, "routing"),
		SearchRouting: stringFromFlags(flags, "search_routing"),
		Filter:        aliasfilter,
	}

	return nil
}
