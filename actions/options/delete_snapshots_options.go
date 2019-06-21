package options

import (
	"encoding/json"
	"errors"

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
