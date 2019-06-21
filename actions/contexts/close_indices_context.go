package contexts

import (
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type closeIndicesContext struct {
	context
}

func (cic *closeIndicesContext) Action() string {
	return "close_indices"
}

func (cic *closeIndicesContext) Config(container *gabs.Container) error {
	return cic.extractConfig(cic.Action(), container, func(container *gabs.Container) error {
		cic.options = new(options.CloseIndicesOptions)

		if err := cic.options.FillFromContainer(container); err != nil {
			return err
		}

		return cic.options.Validate()
	})
}
