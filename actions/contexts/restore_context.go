package contexts

import (
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type restoreContext struct {
	context
}

func (rc *restoreContext) Action() string {
	return "restore"
}

func (rc *restoreContext) Config(container *gabs.Container) error {
	return rc.extractConfig(rc.Action(), container, func(container *gabs.Container) error {
		rc.options = new(options.RestoreOptions)

		if err := rc.options.FillFromContainer(container); err != nil {
			return err
		}

		return rc.options.Validate()
	})
}
