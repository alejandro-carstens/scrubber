package runners

import (
	"scrubber/actions/criterias"
	"scrubber/actions/infos"

	"github.com/alejandro-carstens/golastic"
)

type emptyFilterRunner struct {
	baseRunner
}

// Init initializes the filter runner
func (efr *emptyFilterRunner) Init(connection *golastic.Connection, info ...infos.Informable) (Runnerable, error) {
	err := efr.BaseInit(connection, info...)

	return efr, err
}

// RunFilter filters out elements from the actionable list
func (efr *emptyFilterRunner) RunFilter(channel chan *FilterResponse, criteria criterias.Criteriable) {
	if err := efr.validateCriteria(criteria); err != nil {
		channel <- efr.response.setError(err)
		return
	}

	empty := criteria.(*criterias.Empty)
	passed := false
	docsCount := efr.info.(*infos.IndexInfo).DocsCount

	if docsCount == 0 {
		passed = true
		efr.report.AddReason("Index '%v' is empty", efr.info.Name())
	} else {
		efr.report.AddReason(
			"Index '%v', is not empty, it has '%v' docs",
			efr.info.Name(),
			docsCount,
		)
	}

	channel <- efr.response.setPassed(passed && empty.Include()).setReport(efr.report)
}
