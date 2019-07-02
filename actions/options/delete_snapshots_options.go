package options

import (
	"encoding/json"
	"errors"

	"github.com/spf13/pflag"

	"github.com/Jeffail/gabs"
)

type DeleteSnapshotsOptions struct {
	baseSnapshotOptions
	RetryCount    int `json:"retry_count"`
	RetryInterval int `json:"retry_interval"`
}

func (dso *DeleteSnapshotsOptions) FillFromContainer(container *gabs.Container) error {
	dso.container = container

	return json.Unmarshal(container.Bytes(), dso)
}

func (dso *DeleteSnapshotsOptions) Validate() error {
	if len(dso.Repository) == 0 {
		return errors.New("The repository option is required")
	}

	return nil
}

func (dso *DeleteSnapshotsOptions) BindFlags(flags *pflag.FlagSet) error {
	dso.defaultBindFlags(flags)

	dso.Repository = stringFromFlags(flags, "repository")
	dso.RetryCount = intFromFlags(flags, "retry_count")
	dso.RetryInterval = intFromFlags(flags, "retry_interval")

	return nil
}
