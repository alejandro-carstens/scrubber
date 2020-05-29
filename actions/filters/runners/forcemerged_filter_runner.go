package runners

import (
	"scrubber/actions/criterias"
	"scrubber/actions/infos"

	"github.com/alejandro-carstens/golastic"
)

type forcemergedFilterRunner struct {
	baseRunner
	criteria *criterias.Forcemerged
}

// Init initializes the filter runner
func (ffr *forcemergedFilterRunner) Init(
	criteria criterias.Criteriable,
	connection *golastic.Connection,
	info ...infos.Informable,
) (Runnerable, error) {
	if err := ffr.BaseInit(criteria, connection, info...); err != nil {
		return nil, err
	}

	ffr.criteria = criteria.(*criterias.Forcemerged)

	return ffr, nil
}

// RunFilter filters out elements from the actionable list
func (ffr *forcemergedFilterRunner) RunFilter(channel chan *FilterResponse) {
	segments := ffr.info.(*infos.IndexInfo).SegmentsCount
	passed := segments <= ffr.criteria.MaxNumSegments

	if passed {
		ffr.report.AddReason(
			"Number of segments '%v' is lesser or equal to max number of segments '%v'",
			segments,
			ffr.criteria.MaxNumSegments,
		)
	} else {
		ffr.report.AddReason(
			"Number of segments '%v' is greater than the max number of segments '%v'",
			segments,
			ffr.criteria.MaxNumSegments,
		)
	}

	channel <- &FilterResponse{
		Passed: passed && ffr.criteria.Include(),
		Report: ffr.report,
	}
}
