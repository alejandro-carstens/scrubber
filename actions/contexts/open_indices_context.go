package contexts

import (
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type openIndicesContext struct {
	context
}

func (oic *openIndicesContext) Action() string {
	return "open_indices"
}

func (oic *openIndicesContext) Config(container *gabs.Container) error {
	return oic.extractConfig(oic.Action(), container, func(container *gabs.Container) error {
		oic.options = new(options.OpenIndicesOptions)

		if err := oic.options.FillFromContainer(container); err != nil {
			return err
		}

		return oic.options.Validate()
	})
}
