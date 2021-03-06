package contexts

import (
	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

type DeleteSnapshotsContext struct {
	context
}

func (dsc *DeleteSnapshotsContext) Action() string {
	return "delete_snapshots"
}

func (dsc *DeleteSnapshotsContext) Config(container *gabs.Container) error {
	return dsc.extractConfig(dsc.Action(), container, true, func(container *gabs.Container) error {
		dsc.options = new(options.DeleteSnapshotsOptions)

		if err := dsc.options.FillFromContainer(container); err != nil {
			return err
		}

		return dsc.options.Validate()
	})
}
