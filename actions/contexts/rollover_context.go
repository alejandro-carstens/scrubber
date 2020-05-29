package contexts

import (
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type RolloverContext struct {
	context
}

func (rc *RolloverContext) Action() string {
	return "rollover"
}

func (rc *RolloverContext) Config(container *gabs.Container) error {
	return rc.extractConfig(rc.Action(), container, true, func(container *gabs.Container) error {
		rc.options = new(options.RolloverOptions)

		if err := rc.options.FillFromContainer(container); err != nil {
			return err
		}

		return rc.options.Validate()
	})
}
