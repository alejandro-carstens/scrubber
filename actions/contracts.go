package actions

import (
	"context"

	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/contexts"
	"github.com/alejandro-carstens/scrubber/logger"
	"github.com/alejandro-carstens/scrubber/notifications"
)

// Actionable represents the contract
// to be implemented by actions to be
// performed on an actionable list
// of indices or snapshots
type Actionable interface {
	// Init prepares an action for execution
	Init(
		context contexts.Contextable,
		logger *logger.Logger,
		builder *golastic.Connection,
		queue *notifications.Queue,
		ctx context.Context,
	) error

	// Perform executes the action
	Perform() Actionable

	// ApllyFilter applies different filters
	// through filter runners in order to
	// get an actionable list of indices
	// or snapshots to perform an action
	ApplyFilters() error

	// HasErrors checks if errors occurred
	// while applying the filters to get
	// an actionable list
	HasErrors() bool

	// ApplyOptions applies the action's options
	// prior to performing the given action
	ApplyOptions() Actionable

	// DisableAction prevents the
	// action to be performed
	DisableAction() bool

	// List returns the actionable list
	List() []string

	// Disconnect releases resources used
	// for performing the given action
	Disconnect()

	// Notify issues a notification regarding the execution
	// of an action over an actionable list
	Notify() error
}
