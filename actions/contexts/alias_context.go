package contexts

import (
	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

type AliasContext struct {
	context
}

func (ac *AliasContext) Action() string {
	return "alias"
}

func (ac *AliasContext) Config(container *gabs.Container) error {
	return ac.extractConfig(ac.Action(), container, true, func(container *gabs.Container) error {
		ac.options = new(options.AliasOptions)

		if err := ac.options.FillFromContainer(container); err != nil {
			return err
		}

		return ac.options.Validate()
	})
}
