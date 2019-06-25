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

	if filter, _ := flags.GetString("filter"); len(filter) > 0 {
		if err := json.Unmarshal([]byte(filter), &aliasfilter); err != nil {
			return err
		}
	}

	name, _ := flags.GetString("name")
	aliasType, _ := flags.GetString("type")
	routing, _ := flags.GetString("routing")
	searchRouting, _ := flags.GetString("search_routing")

	ao.Name = name
	ao.Type = aliasType
	ao.ExtraSettings = &AliasExtraSettingsOption{
		Routing:       routing,
		SearchRouting: searchRouting,
		Filter:        aliasfilter,
	}

	return nil
}
