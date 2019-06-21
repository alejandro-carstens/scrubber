package contexts

import (
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type indexSettingsContext struct {
	context
}

func (isc *indexSettingsContext) Action() string {
	return "index_settings"
}

func (isc *indexSettingsContext) Config(container *gabs.Container) error {
	return isc.extractConfig(isc.Action(), container, func(container *gabs.Container) error {
		isc.options = new(options.IndexSettingsOptions)

		if err := isc.options.FillFromContainer(container); err != nil {
			return err
		}

		return isc.options.Validate()
	})
}
