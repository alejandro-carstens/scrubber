package runners

import (
	"errors"

	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/criterias"
	"github.com/alejandro-carstens/scrubber/actions/infos"
)

const CLOSED_STATUS string = "close"
const OPEN_STATUS string = "open"

type closedFilterRunner struct {
	baseRunner
	criteria *criterias.Closed
}

// Init initializes the filter runner
func (cfr *closedFilterRunner) Init(
	criteria criterias.Criteriable,
	connection *golastic.Connection,
	info ...infos.Informable,
) (Runnerable, error) {
	if err := cfr.BaseInit(criteria, connection, info...); err != nil {
		return nil, err
	}

	cfr.criteria = criteria.(*criterias.Closed)

	return cfr, nil
}

// RunFilter filters out elements from the actionable list
func (cfr *closedFilterRunner) RunFilter(channel chan *FilterResponse) {
	var err error

	passed := false
	status := cfr.info.(*infos.IndexInfo).Status

	if len(status) == 0 {
		err = errors.New("Could not determine if the index is closed")

		passed = false
	}

	if status == CLOSED_STATUS {
		cfr.report.AddReason("Index '%v' is closed", cfr.info.Name())

		passed = true
	}

	if !passed {
		cfr.report.AddReason("Index '%v' is  not closed", cfr.info.Name())
	}

	channel <- &FilterResponse{
		Err:    err,
		Passed: passed && cfr.criteria.Include(),
		Report: cfr.report,
	}
}
