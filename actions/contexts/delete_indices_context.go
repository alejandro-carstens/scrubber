package contexts

import (
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type deleteIndicesContext struct {
	context
}

func (dic *deleteIndicesContext) Action() string {
	return "delete_indices"
}

func (dic *deleteIndicesContext) Config(container *gabs.Container) error {
	return dic.extractConfig(dic.Action(), container, func(container *gabs.Container) error {
		dic.options = new(options.DeleteIndicesOptions)

		if err := dic.options.FillFromContainer(container); err != nil {
			return err
		}

		return dic.options.Validate()
	})
}
