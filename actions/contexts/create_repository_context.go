package contexts

import (
	"errors"
	"fmt"
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type createRepositoryContext struct {
	context
}

func (crc *createRepositoryContext) Action() string {
	return "create_repository"
}

func (crc *createRepositoryContext) Config(container *gabs.Container) error {
	config, err := container.ChildrenMap()

	if err != nil {
		return err
	}

	if len(config) == 0 {
		return errors.New("Config is empty")
	}

	crc.config = container

	value, valid := config["action"]

	if !valid || fmt.Sprint(value.Data()) != crc.Action() {
		return errors.New("action not of type create_repository")
	}

	options, valid := config["options"]

	if !valid {
		return crc.marshallOptions(nil)
	}

	return crc.marshallOptions(options)
}

func (crc *createRepositoryContext) marshallOptions(container *gabs.Container) error {
	crc.options = new(options.CreateRepositoryOptions)

	if err := crc.options.FillFromContainer(container); err != nil {
		return err
	}

	return crc.options.Validate()
}
