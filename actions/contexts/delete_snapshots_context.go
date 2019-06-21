package contexts

import (
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type deleteSnapshotsContext struct {
	context
}

func (dsc *deleteSnapshotsContext) Action() string {
	return "delete_snapshots"
}

func (dsc *deleteSnapshotsContext) Config(container *gabs.Container) error {
	return dsc.extractConfig(dsc.Action(), container, func(container *gabs.Container) error {
		dsc.options = new(options.DeleteSnapshotsOptions)

		if err := dsc.options.FillFromContainer(container); err != nil {
			return err
		}

		return dsc.options.Validate()
	})
}
