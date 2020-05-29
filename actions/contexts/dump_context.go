package contexts

import (
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type DumpContext struct {
	context
}

func (dc *DumpContext) Action() string {
	return "dump"
}

func (dc *DumpContext) Config(container *gabs.Container) error {
	return dc.extractConfig(dc.Action(), container, true, func(container *gabs.Container) error {
		dc.options = new(options.DumpOptions)

		if err := dc.options.FillFromContainer(container); err != nil {
			return err
		}

		return dc.options.Validate()
	})
}
