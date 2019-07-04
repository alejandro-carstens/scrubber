package contexts

import (
	"errors"
	"fmt"
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type DeleteRepositoriesContext struct {
	context
}

func (drc *DeleteRepositoriesContext) Action() string {
	return "delete_repositories"
}

func (drc *DeleteRepositoriesContext) Config(container *gabs.Container) error {
	config, err := container.ChildrenMap()

	if err != nil {
		return err
	}

	if len(config) == 0 {
		return errors.New("Config is empty")
	}

	drc.config = container

	value, valid := config["action"]

	if !valid || fmt.Sprint(value.Data()) != drc.Action() {
		return errors.New("action not of type create_repository")
	}

	options, valid := config["options"]

	if !valid {
		return drc.marshallOptions(nil)
	}

	return drc.marshallOptions(options)
}

func (drc *DeleteRepositoriesContext) marshallOptions(container *gabs.Container) error {
	drc.options = new(options.DeleteRepositoriesOptions)

	if err := drc.options.FillFromContainer(container); err != nil {
		return err
	}

	return drc.options.Validate()
}
