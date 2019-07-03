package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type DeleteRepositoriesOptions struct {
	defaultOptions
	Repositories string `json:"repositories"`
}

func (dro *DeleteRepositoriesOptions) FillFromContainer(container *gabs.Container) error {
	dro.container = container

	return json.Unmarshal(container.Bytes(), dro)
}

func (dro *DeleteRepositoriesOptions) Validate() error {
	if len(dro.Repositories) == 0 {
		return errors.New("repositories is a required field")
	}

	return nil
}

func (dro *DeleteRepositoriesOptions) BindFlags(flags *pflag.FlagSet) error {
	dro.defaultBindFlags(flags)

	dro.Repositories = stringFromFlags(flags, "repositories")

	return nil
}
