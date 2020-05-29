package contexts

import (
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type MutateContext struct {
	context
}

func (mc *MutateContext) Action() string {
	return "mutate"
}

func (mc *MutateContext) Config(container *gabs.Container) error {
	return mc.extractConfig(mc.Action(), container, true, func(container *gabs.Container) error {
		mc.options = new(options.MutateOptions)

		if err := mc.options.FillFromContainer(container); err != nil {
			return err
		}

		return mc.options.Validate()
	})
}
