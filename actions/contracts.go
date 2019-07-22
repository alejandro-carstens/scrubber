package actions

import (
	"scrubber/actions/contexts"
	"scrubber/logger"
)

type Actionable interface {
	Init(context contexts.Contextable, logger *logger.Logger) error

	Perform() Actionable

	ApplyFilters() error

	HasErrors() bool

	ApplyOptions() Actionable

	DisableAction() bool

	List() []string

	TearDownBuilder()
}
