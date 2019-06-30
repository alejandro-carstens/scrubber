package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type ListSnapshotsOptions struct {
	baseSnapshotOptions
}

func (lso *ListSnapshotsOptions) FillFromContainer(container *gabs.Container) error {
	lso.container = container

	return json.Unmarshal(container.Bytes(), lso)
}

func (lso *ListSnapshotsOptions) Validate() error {
	if len(lso.Repository) == 0 {
		return errors.New("repository value is required")
	}

	return nil
}

func (lso *ListSnapshotsOptions) BindFlags(flags *pflag.FlagSet) error {
	lso.defaultBindFlags(flags)

	repository, _ := flags.GetString("repository")

	lso.Repository = repository

	return nil
}
