package contexts

import (
	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

type DeleteRepositoriesContext struct {
	context
}

func (drc *DeleteRepositoriesContext) Action() string {
	return "delete_repositories"
}

func (drc *DeleteRepositoriesContext) Config(container *gabs.Container) error {
	return drc.extractConfig(drc.Action(), container, false, func(container *gabs.Container) error {
		drc.options = new(options.DeleteRepositoriesOptions)

		if err := drc.options.FillFromContainer(container); err != nil {
			return err
		}

		return drc.options.Validate()
	})
}
