package criterias

import (
	"github.com/Jeffail/gabs"
)

type Criteriable interface {
	Validate() error

	FillFromContainer(container *gabs.Container) (Criteriable, error)

	Name() string
}

type Sortable interface {
	GetField() string

	GetStatsResult() string

	GetReverse() bool

	GetTimestring() string

	GetSource() string

	GetStrictMode() bool
}
