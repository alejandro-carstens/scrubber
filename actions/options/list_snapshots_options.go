package options

import (
	"encoding/json"
	"errors"
	"fmt"

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
	repository, _ := flags.GetString("repository")

	lso.Repository = repository

	return lso.defaultBindFlags(flags)
}

func (lso *ListSnapshotsOptions) Exists(value string) bool {
	if lso.container == nil {
		lso.container = toContainer(lso)
	}

	return lso.container.Exists(value)
}

func (lso *ListSnapshotsOptions) String(value string) string {
	if lso.container == nil {
		lso.container = toContainer(lso)
	}

	return fmt.Sprint(lso.container.S(value).Data())
}
