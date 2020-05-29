package runners

import (
	"scrubber/actions/criterias"
	"scrubber/actions/infos"

	"github.com/alejandro-carstens/golastic"
)

type allocatedFilterRunner struct {
	baseRunner
	criteria *criterias.Allocated
}

// Init initializes the filter runner
func (afr *allocatedFilterRunner) Init(
	criteria criterias.Criteriable,
	connection *golastic.Connection,
	info ...infos.Informable,
) (Runnerable, error) {
	if err := afr.BaseInit(criteria, connection, info...); err != nil {
		return nil, err
	}

	afr.criteria = criteria.(*criterias.Allocated)

	return afr, nil
}

// RunFilter filters out elements from the actionable list
func (afr *allocatedFilterRunner) RunFilter(channel chan *FilterResponse) {
	settingsResponse, err := afr.connection.Indexer(nil).Settings(afr.info.Name())

	if err != nil {
		channel <- &FilterResponse{Err: err}

		return
	}

	passed := false
	container := settingsResponse[afr.info.Name()]

	value, valid := container.S("index", "routing", "allocation", afr.criteria.AllocationType, afr.criteria.Key).Data().(string)

	if valid {
		passed = value == afr.criteria.Value

		if passed {
			afr.report.AddReason(
				"Allocation value '%v' matched for allocation type '%v' and key '%v'",
				value,
				afr.criteria.AllocationType,
				afr.criteria.Key,
			)
		}
	}

	if !passed {
		afr.report.AddReason(
			"Allocation value '%v' not matched for allocation type '%v' and key '%v'",
			value,
			afr.criteria.AllocationType,
			afr.criteria.Key,
		)
	}

	channel <- &FilterResponse{
		Err:    err,
		Passed: passed && afr.criteria.Include(),
		Report: afr.report,
	}
}
