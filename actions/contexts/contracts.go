package contexts

import (
	"scrubber/actions/criterias"
	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
)

type Contextable interface {
	Action() string

	Container() *gabs.Container

	GetAsync() bool

	GetNumberOfWorkers() int

	GetQueueLength() int

	GetRetryCount() int

	Options() options.Optionable

	Builder() *criterias.Builder

	Config(container *gabs.Container) error

	setNumberOfWorkers(numberOfWorkers int)

	setQueueLength(queueLength int)
}
