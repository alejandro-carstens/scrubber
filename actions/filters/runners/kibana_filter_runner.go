package runners

import (
	"scrubber/actions/criterias"
	"scrubber/actions/infos"
	"strings"

	"github.com/alejandro-carstens/golastic"
)

type kibanaFilterRunner struct {
	baseRunner
}

// Init initializes the filter runner
func (kfr *kibanaFilterRunner) Init(connection *golastic.Connection, info ...infos.Informable) (Runnerable, error) {
	err := kfr.BaseInit(connection, info...)

	return kfr, err
}

// RunFilter filters out elements from the actionable list
func (kfr *kibanaFilterRunner) RunFilter(channel chan *FilterResponse, criteria criterias.Criteriable) {
	if err := kfr.validateCriteria(criteria); err != nil {
		channel <- kfr.response.setError(err)
		return
	}

	isKibana := strings.HasPrefix(kfr.info.Name(), ".kibana")

	if isKibana {
		kfr.report.AddReason("Index '%v' is a kibana index", kfr.info.Name())
	} else {
		kfr.report.AddReason("Index '%v' is  not a kibana index", kfr.info.Name())
	}

	kibana := criteria.(*criterias.Kibana)

	channel <- kfr.response.setPassed(isKibana && kibana.Include()).setReport(kfr.report)
}
