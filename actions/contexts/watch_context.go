package contexts

import (
	"github.com/alejandro-carstens/scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type WatchContext struct {
	context
}

func (wc *WatchContext) Action() string {
	return "watch"
}

func (wc *WatchContext) Config(container *gabs.Container) error {
	return wc.extractConfig(wc.Action(), container, true, func(container *gabs.Container) error {
		wc.options = new(options.WatchOptions)

		if err := wc.options.FillFromContainer(container); err != nil {
			return err
		}

		return wc.options.Validate()
	})
}
