package contexts

import (
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type AliasContext struct {
	context
}

func (ac *AliasContext) Action() string {
	return "alias"
}

func (ac *AliasContext) Config(container *gabs.Container) error {
	return ac.extractConfig(ac.Action(), container, func(container *gabs.Container) error {
		ac.options = new(options.AliasOptions)

		if err := ac.options.FillFromContainer(container); err != nil {
			return err
		}

		return ac.options.Validate()
	})
}
