package actions

import (
	"scrubber/actions/contexts"
	"scrubber/logger"

	"github.com/alejandro-carstens/golastic"
)

type Actionable interface {
	Init(context contexts.Contextable, logger *logger.Logger, builder *golastic.Connection) error

	Perform() Actionable

	ApplyFilters() error

	HasErrors() bool

	ApplyOptions() Actionable

	DisableAction() bool

	List() []string

	Disconnect()
}
