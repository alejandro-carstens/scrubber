package contexts

import (
	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

type ListSnapshotsContext struct {
	context
}

func (ls *ListSnapshotsContext) Action() string {
	return "list_snapshots"
}

func (ls *ListSnapshotsContext) Config(container *gabs.Container) error {
	return ls.extractConfig(ls.Action(), container, true, func(container *gabs.Container) error {
		ls.options = new(options.ListSnapshotsOptions)

		if err := ls.options.FillFromContainer(container); err != nil {
			return err
		}

		return ls.options.Validate()
	})
}
