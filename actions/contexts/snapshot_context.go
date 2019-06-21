package contexts

import (
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type snapshotContext struct {
	context
}

func (sc *snapshotContext) Action() string {
	return "snapshot"
}

func (sc *snapshotContext) Config(container *gabs.Container) error {
	return sc.extractConfig(sc.Action(), container, func(container *gabs.Container) error {
		sc.options = new(options.SnapshotOptions)

		if err := sc.options.FillFromContainer(container); err != nil {
			return err
		}

		return sc.options.Validate()
	})
}
