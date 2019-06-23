package reports

import (
	"scrubber/actions/criterias"

	"github.com/Jeffail/gabs"
)

type Reportable interface {
	Line() (string, error)

	// SetCriteria sets the criteria used by the filter runner
	SetCriteria(criteria criterias.Criteriable)

	// ToJson converts the report to JSON
	ToJson() (*gabs.Container, error)

	AddReason(reason string, values ...interface{})

	// Error adds an error message to the summary
	Error(err error) Reportable
}
