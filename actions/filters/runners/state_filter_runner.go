package runners

import (
	"scrubber/actions/criterias"
	"scrubber/actions/infos"

	"github.com/alejandro-carstens/golastic"
)

type stateFilterRunner struct {
	baseRunner
}

// Init initializes the filter runner
func (sfr *stateFilterRunner) Init(connection *golastic.Connection, info ...infos.Informable) (Runnerable, error) {
	err := sfr.BaseInit(connection, info...)

	return sfr, err
}

// RunFilter filters out elements from the actionable list
func (sfr *stateFilterRunner) RunFilter(channel chan *FilterResponse, criteria criterias.Criteriable) {
	if err := sfr.validateCriteria(criteria); err != nil {
		channel <- sfr.response.setError(err)
		return
	}

	state := criteria.(*criterias.State)
	snapshotInfo := sfr.info.(*infos.SnapshotInfo)

	passed := state.State == snapshotInfo.State

	if passed {
		sfr.report.AddReason(
			"Snapshot '%v' state '%v' matches '%v'",
			sfr.info.Name(),
			state.State,
			snapshotInfo.State,
		)
	} else {
		sfr.report.AddReason(
			"Snapshot '%v' state '%v' does not match '%v'",
			sfr.info.Name(),
			state.State,
			snapshotInfo.State,
		)
	}

	channel <- sfr.response.setReport(sfr.report).setPassed(passed && state.Include())
}
