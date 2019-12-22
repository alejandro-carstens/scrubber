package contexts

import (
	"github.com/alejandro-carstens/scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type CreateRepositoryContext struct {
	context
}

func (crc *CreateRepositoryContext) Action() string {
	return "create_repository"
}

func (crc *CreateRepositoryContext) Config(container *gabs.Container) error {
	return crc.extractConfig(crc.Action(), container, false, func(container *gabs.Container) error {
		crc.options = new(options.CreateRepositoryOptions)

		if err := crc.options.FillFromContainer(container); err != nil {
			return err
		}

		return crc.options.Validate()
	})
}
