package contexts

import (
	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

type OpenIndicesContext struct {
	context
}

func (oic *OpenIndicesContext) Action() string {
	return "open_indices"
}

func (oic *OpenIndicesContext) Config(container *gabs.Container) error {
	return oic.extractConfig(oic.Action(), container, true, func(container *gabs.Container) error {
		oic.options = new(options.OpenIndicesOptions)

		if err := oic.options.FillFromContainer(container); err != nil {
			return err
		}

		return oic.options.Validate()
	})
}
