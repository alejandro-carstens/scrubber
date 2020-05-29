package contexts

import (
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type CloseIndicesContext struct {
	context
}

func (cic *CloseIndicesContext) Action() string {
	return "close_indices"
}

func (cic *CloseIndicesContext) Config(container *gabs.Container) error {
	return cic.extractConfig(cic.Action(), container, true, func(container *gabs.Container) error {
		cic.options = new(options.CloseIndicesOptions)

		if err := cic.options.FillFromContainer(container); err != nil {
			return err
		}

		return cic.options.Validate()
	})
}
