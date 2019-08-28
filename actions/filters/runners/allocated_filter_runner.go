package runners

import (
	"scrubber/actions/criterias"
	"scrubber/actions/infos"

	"github.com/alejandro-carstens/golastic"
)

type allocatedFilterRunner struct {
	baseRunner
}

// Init initializes the filter runner
func (afr *allocatedFilterRunner) Init(connection *golastic.Connection, info ...infos.Informable) (Runnerable, error) {
	if err := afr.BaseInit(connection, info...); err != nil {
		return nil, err
	}

	return afr, nil
}

// RunFilter filters out elements from the actionable list
func (afr *allocatedFilterRunner) RunFilter(channel chan *FilterResponse, criteria criterias.Criteriable) {
	if err := afr.validateCriteria(criteria); err != nil {
		channel <- afr.response.setError(err)
		return
	}

	allocated := criteria.(*criterias.Allocated)

	settingsResponse, err := afr.connection.Indexer(nil).Settings(afr.info.Name())

	if err != nil {
		channel <- afr.response.setError(err)
		return
	}

	passed := false
	container := settingsResponse[afr.info.Name()]

	value, valid := container.S("index", "routing", "allocation", allocated.AllocationType, allocated.Key).Data().(string)

	if valid {
		passed = value == allocated.Value

		if passed {
			afr.report.AddReason(
				"Allocation value '%v' matched for allocation type '%v' and key '%v'",
				value,
				allocated.AllocationType,
				allocated.Key,
			)
		}
	}

	if !passed {
		afr.report.AddReason(
			"Allocation value '%v' not matched for allocation type '%v' and key '%v'",
			value,
			allocated.AllocationType,
			allocated.Key,
		)
	}

	channel <- afr.response.setPassed(passed && allocated.Include()).setReport(afr.report)
}
