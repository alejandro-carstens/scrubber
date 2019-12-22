package contexts

import (
	"github.com/alejandro-carstens/scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type SnapshotContext struct {
	context
}

func (sc *SnapshotContext) Action() string {
	return "snapshot"
}

func (sc *SnapshotContext) Config(container *gabs.Container) error {
	return sc.extractConfig(sc.Action(), container, true, func(container *gabs.Container) error {
		sc.options = new(options.SnapshotOptions)

		if err := sc.options.FillFromContainer(container); err != nil {
			return err
		}

		return sc.options.Validate()
	})
}
