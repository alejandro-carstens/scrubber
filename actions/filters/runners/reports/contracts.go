package reports

import (
	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/scrubber/actions/criterias"
)

// Reportable is a contract implemented by structs that report
// on the activity of a filter runner
type Reportable interface {
	// Line returns a human readable string or a filter runner activity
	Line() (string, error)

	// SetCriteria sets the criteria used by the filter runner
	SetCriteria(criteria criterias.Criteriable)

	// ToJson converts the report to JSON
	ToJson() (*gabs.Container, error)

	// AddReason appends a reason for why an element of
	// the actionable list remained or was excluded
	AddReason(reason string, values ...interface{})

	// Error adds an error message to the summary
	Error(err error) Reportable
}
