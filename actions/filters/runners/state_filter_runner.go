package runners

import (
	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/criterias"
	"github.com/alejandro-carstens/scrubber/actions/infos"
)

type stateFilterRunner struct {
	baseRunner
	criteria *criterias.State
}

// Init initializes the filter runner
func (sfr *stateFilterRunner) Init(criteria criterias.Criteriable, connection *golastic.Connection, info ...infos.Informable) (Runnerable, error) {
	if err := sfr.BaseInit(criteria, connection, info...); err != nil {
		return nil, err
	}

	sfr.criteria = criteria.(*criterias.State)

	return sfr, nil
}

// RunFilter filters out elements from the actionable list
func (sfr *stateFilterRunner) RunFilter(channel chan *FilterResponse) {
	snapshotInfo := sfr.info.(*infos.SnapshotInfo)

	passed := sfr.criteria.State == snapshotInfo.State

	if passed {
		sfr.report.AddReason(
			"Snapshot '%v' state '%v' matches '%v'",
			sfr.info.Name(),
			sfr.criteria.State,
			snapshotInfo.State,
		)
	} else {
		sfr.report.AddReason(
			"Snapshot '%v' state '%v' does not match '%v'",
			sfr.info.Name(),
			sfr.criteria.State,
			snapshotInfo.State,
		)
	}

	channel <- sfr.response.setReport(sfr.report).setPassed(passed && sfr.criteria.Include())
}
