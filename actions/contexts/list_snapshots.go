package contexts

import (
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type ListSnapshots struct {
	context
}

func (ls *ListSnapshots) Action() string {
	return "list_snapshots"
}

func (ls *ListSnapshots) Config(container *gabs.Container) error {
	return ls.extractConfig(ls.Action(), container, func(container *gabs.Container) error {
		ls.options = new(options.ListSnapshotsOptions)

		if err := ls.options.FillFromContainer(container); err != nil {
			return err
		}

		return ls.options.Validate()
	})
}
