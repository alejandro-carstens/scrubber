package contexts

import (
	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

type RestoreContext struct {
	context
}

func (rc *RestoreContext) Action() string {
	return "restore"
}

func (rc *RestoreContext) Config(container *gabs.Container) error {
	return rc.extractConfig(rc.Action(), container, true, func(container *gabs.Container) error {
		rc.options = new(options.RestoreOptions)

		if err := rc.options.FillFromContainer(container); err != nil {
			return err
		}

		return rc.options.Validate()
	})
}
