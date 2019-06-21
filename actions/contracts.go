package actions

import (
	"scrubber/actions/contexts"
	"scrubber/logging"
)

type Actionable interface {
	Init(context contexts.Contextable, logger *logging.SrvLogger) error

	Perform() Actionable

	ApplyFilters() error

	HasErrors() bool

	ApplyOptions() Actionable

	DisableAction() bool
}
