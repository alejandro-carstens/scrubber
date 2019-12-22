package contexts

import (
	"github.com/alejandro-carstens/scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type DeleteIndicesContext struct {
	context
}

func (dic *DeleteIndicesContext) Action() string {
	return "delete_indices"
}

func (dic *DeleteIndicesContext) Config(container *gabs.Container) error {
	return dic.extractConfig(dic.Action(), container, true, func(container *gabs.Container) error {
		dic.options = new(options.DeleteIndicesOptions)

		if err := dic.options.FillFromContainer(container); err != nil {
			return err
		}

		return dic.options.Validate()
	})
}
