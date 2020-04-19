package contexts

import (
	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

type ImportDumpContext struct {
	context
}

func (idc *ImportDumpContext) Action() string {
	return "import_dump"
}

func (idc *ImportDumpContext) Config(container *gabs.Container) error {
	return idc.extractConfig(idc.Action(), container, true, func(container *gabs.Container) error {
		idc.options = new(options.ImportDumpOptions)

		if err := idc.options.FillFromContainer(container); err != nil {
			return err
		}

		return idc.options.Validate()
	})
}
