package options

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type ListIndicesOptions struct {
	defaultOptions
}

func (lio *ListIndicesOptions) FillFromContainer(container *gabs.Container) error {
	lio.container = container

	return json.Unmarshal(container.Bytes(), lio)
}

func (lo *ListIndicesOptions) Validate() error {
	return nil
}

func (lio *ListIndicesOptions) BindFlags(flags *pflag.FlagSet) error {
	return lio.defaultBindFlags(flags)
}
