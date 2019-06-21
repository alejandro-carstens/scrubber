package contexts

import (
	"errors"
	"fmt"
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type createIndexContext struct {
	context
}

func (cic *createIndexContext) Action() string {
	return "create_index"
}

func (cic *createIndexContext) Config(container *gabs.Container) error {
	config, err := container.ChildrenMap()

	if err != nil {
		return err
	}

	if len(config) == 0 {
		return errors.New("Config is empty")
	}

	cic.config = container

	value, valid := config["action"]

	if !valid || fmt.Sprint(value.Data()) != cic.Action() {
		return errors.New("action not of type create_index")
	}

	options, valid := config["options"]

	if !valid {
		return cic.marshallOptions(nil)
	}

	return cic.marshallOptions(options)
}

func (cic *createIndexContext) marshallOptions(container *gabs.Container) error {
	cic.options = new(options.CreateIndexOptions)

	if err := cic.options.FillFromContainer(container); err != nil {
		return err
	}

	return cic.options.Validate()
}
