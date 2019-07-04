package contexts

import (
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type CreateIndexContext struct {
	context
}

func (cic *CreateIndexContext) Action() string {
	return "create_index"
}

func (cic *CreateIndexContext) Config(container *gabs.Container) error {
	return cic.extractConfig(cic.Action(), container, false, func(container *gabs.Container) error {
		cic.options = new(options.CreateIndexOptions)

		if err := cic.options.FillFromContainer(container); err != nil {
			return err
		}

		return cic.options.Validate()
	})
}
