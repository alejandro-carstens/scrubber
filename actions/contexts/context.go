package contexts

import (
	"errors"
	"fmt"
	"scrubber/actions/criterias"
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type extractOptions func(options *gabs.Container) error

type context struct {
	config          *gabs.Container
	options         options.Optionable
	builder         *criterias.Builder
	List            []string `json:"list"`
	Async           bool     `json:"async"`
	NumberOfWorkers int      `json:"number_of_workers"`
	QueueLength     int      `json:"queue_length"`
	RetryCount      int      `json:"retry_count"`
}

func (c *context) Builder() *criterias.Builder {
	if c.builder == nil {
		c.builder = criterias.NewBuilder()
	}

	return c.builder
}

func (c *context) Container() *gabs.Container {
	return c.config
}

func (c *context) GetAsync() bool {
	return c.Async
}

func (c *context) GetNumberOfWorkers() int {
	return c.NumberOfWorkers
}

func (c *context) GetQueueLength() int {
	return c.QueueLength
}

func (c *context) GetRetryCount() int {
	return c.RetryCount
}

func (c *context) Options() options.Optionable {
	return c.options
}

func (c *context) SetOptions(options options.Optionable) error {
	c.options = options

	return c.options.Validate()
}

func (c *context) SetList(list ...string) {
	c.List = list
}

func (c *context) ActionableList() []string {
	return c.List
}

func (c *context) setNumberOfWorkers(numberOfWorkers int) {
	c.NumberOfWorkers = numberOfWorkers
}

func (c *context) setQueueLength(queueLength int) {
	c.QueueLength = queueLength
}

func (c *context) extractConfig(action string, container *gabs.Container, extractFilters bool, fn extractOptions) error {
	config, err := container.ChildrenMap()

	if err != nil {
		return err
	}

	if len(config) == 0 {
		return errors.New("Config is empty")
	}

	c.config = container

	if value, valid := config["action"]; !valid || fmt.Sprint(value.Data()) != action {
		return errors.New("action not of type " + action)
	}

	if extractFilters {
		if err := c.extractFilters(action, config); err != nil {
			return err
		}
	}

	options, valid := config["options"]

	if !valid {
		return fn(nil)
	}

	return fn(options)
}

func (c *context) extractFilters(action string, config map[string]*gabs.Container) error {
	filterContainer, valid := config["filters"]

	if !valid {
		return errors.New("Filters are required")
	}

	filters, err := filterContainer.Children()

	if err != nil {
		return err
	}

	if isSnapshotAction(action) {
		if err := c.validateSnapshotAction(filters); err != nil {
			return err
		}
	}

	c.builder = criterias.NewBuilder()

	return c.builder.Build(action, filters)
}

func (c *context) validateSnapshotAction(filters []*gabs.Container) error {
	for _, filter := range filters {
		filterType, valid := filter.S("filtertype").Data().(string)

		if !valid {
			return errors.New("Could not retrieve filtertype value")
		}

		if filterType != "count" && filterType != "age" && filterType != "pattern" && filterType != "state" {
			return errors.New("Invalid filtertype for snapshot action")
		}
	}

	return nil
}
