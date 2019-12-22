package contexts

import (
	"github.com/alejandro-carstens/scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type ListIndicesContext struct {
	context
}

func (lic *ListIndicesContext) Action() string {
	return "list_indices"
}

func (lic *ListIndicesContext) Config(container *gabs.Container) error {
	return lic.extractConfig(lic.Action(), container, true, func(container *gabs.Container) error {
		lic.options = new(options.ListIndicesOptions)

		if err := lic.options.FillFromContainer(container); err != nil {
			return err
		}

		return lic.options.Validate()
	})
}
