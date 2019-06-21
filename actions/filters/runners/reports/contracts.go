package reports

import (
	"scrubber/actions/criterias"

	"github.com/Jeffail/gabs"
)

type Reportable interface {
	Line() (string, error)

	SetCriteria(criteria criterias.Criteriable)

	ToJson() (*gabs.Container, error)

	AddReason(reason string, values ...interface{})

	Error(err error) Reportable
}
