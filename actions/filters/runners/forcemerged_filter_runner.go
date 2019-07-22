package runners

import (
	"scrubber/actions/criterias"
	"scrubber/actions/infos"

	"github.com/alejandro-carstens/golastic"
)

type forcemergedFilterRunner struct {
	baseRunner
	builder golastic.Queryable
}

// Init initializes the filter runner
func (ffr *forcemergedFilterRunner) Init(builder *golastic.ElasticsearchBuilder, info ...infos.Informable) (Runnerable, error) {
	if err := ffr.BaseInit(builder, info...); err != nil {
		return nil, err
	}

	return ffr, nil
}

// RunFilter filters out elements from the actionable list
func (ffr *forcemergedFilterRunner) RunFilter(channel chan *FilterResponse, criteria criterias.Criteriable) {
	if err := ffr.validateCriteria(criteria); err != nil {
		channel <- ffr.response.setError(err)
		return
	}

	segments := ffr.info.(*infos.IndexInfo).SegmentsCount

	forcemerged := criteria.(*criterias.Forcemerged)

	passed := segments <= forcemerged.MaxNumSegments

	if passed {
		ffr.report.AddReason(
			"Number of segments '%v' is lesser or equal to max number of segments '%v'",
			segments,
			forcemerged.MaxNumSegments,
		)
	} else {
		ffr.report.AddReason(
			"Number of segments '%v' is greater than the max number of segments '%v'",
			segments,
			forcemerged.MaxNumSegments,
		)
	}

	channel <- ffr.response.setPassed(passed && forcemerged.Include()).setReport(ffr.report)
}
