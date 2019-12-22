package runners

import (
	"github.com/alejandro-carstens/scrubber/actions/criterias"
	"github.com/alejandro-carstens/scrubber/actions/infos"

	"github.com/alejandro-carstens/golastic"
)

type emptyFilterRunner struct {
	baseRunner
	criteria *criterias.Empty
}

// Init initializes the filter runner
func (efr *emptyFilterRunner) Init(criteria criterias.Criteriable, connection *golastic.Connection, info ...infos.Informable) (Runnerable, error) {
	if err := efr.BaseInit(criteria, connection, info...); err != nil {
		return nil, err
	}

	efr.criteria = criteria.(*criterias.Empty)

	return efr, nil
}

// RunFilter filters out elements from the actionable list
func (efr *emptyFilterRunner) RunFilter(channel chan *FilterResponse) {
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

	channel <- efr.response.setPassed(passed && efr.criteria.Include()).setReport(efr.report)
}
