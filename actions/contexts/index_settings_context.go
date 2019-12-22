package contexts

import (
	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

type IndexSettingsContext struct {
	context
}

func (isc *IndexSettingsContext) Action() string {
	return "index_settings"
}

func (isc *IndexSettingsContext) Config(container *gabs.Container) error {
	return isc.extractConfig(isc.Action(), container, true, func(container *gabs.Container) error {
		isc.options = new(options.IndexSettingsOptions)

		if err := isc.options.FillFromContainer(container); err != nil {
			return err
		}

		return isc.options.Validate()
	})
}
