package contexts

import (
	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

type ImportContext struct {
	context
}

func (ic *ImportContext) Action() string {
	return "import"
}

func (ic *ImportContext) Config(container *gabs.Container) error {
	return ic.extractConfig(ic.Action(), container, true, func(container *gabs.Container) error {
		ic.options = new(options.ImportOptions)

		if err := ic.options.FillFromContainer(container); err != nil {
			return err
		}

		return ic.options.Validate()
	})
}
