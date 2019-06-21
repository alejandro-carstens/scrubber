package runners

import (
	"errors"
	"scrubber/actions/criterias"
	"scrubber/actions/responses"
)

const CLOSED_STATUS string = "close"
const OPEN_STATUS string = "open"

type closedFilterRunner struct {
	baseRunner
}

func (cfr *closedFilterRunner) Init(info ...responses.Informable) (Runnerable, error) {
	if err := cfr.BaseInit(info...); err != nil {
		return nil, err
	}

	return cfr, nil
}

func (cfr *closedFilterRunner) RunFilter(channel chan *FilterResponse, criteria criterias.Criteriable) {
	if err := cfr.validateCriteria(criteria); err != nil {
		channel <- cfr.response.setError(err)
		return
	}

	closed := criteria.(*criterias.Closed)

	passed := false

	status := cfr.info.(*responses.IndexInfo).Status

	if len(status) == 0 {
		cfr.response.setError(errors.New("Could not determine if the index is closed"))
		passed = false
	}

	if status == CLOSED_STATUS {
		cfr.report.AddReason("Index '%v' is closed", cfr.info.Name())
		passed = true
	}

	if !passed {
		cfr.report.AddReason("Index '%v' is  not closed", cfr.info.Name())
	}

	channel <- cfr.response.setPassed(passed && closed.Include()).setReport(cfr.report)
}
