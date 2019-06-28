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

	repository, _ := flags.GetString("repository")
	retryCount, _ := flags.GetInt("retry_count")
	retryInterval, _ := flags.GetInt("retry_interval")

	dso.Repository = repository
	dso.RetryCount = retryCount
	dso.RetryInterval = retryInterval

	return nil
}
