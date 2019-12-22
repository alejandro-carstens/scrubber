package runners

import (
	"github.com/alejandro-carstens/scrubber/actions/criterias"
	"github.com/alejandro-carstens/scrubber/actions/infos"
	"strings"

	"github.com/alejandro-carstens/golastic"
)

type kibanaFilterRunner struct {
	baseRunner
	criteria *criterias.Kibana
}

// Init initializes the filter runner
func (kfr *kibanaFilterRunner) Init(criteria criterias.Criteriable, connection *golastic.Connection, info ...infos.Informable) (Runnerable, error) {
	if err := kfr.BaseInit(criteria, connection, info...); err != nil {
		return nil, err
	}

	kfr.criteria = criteria.(*criterias.Kibana)

	return kfr, nil
}

// RunFilter filters out elements from the actionable list
func (kfr *kibanaFilterRunner) RunFilter(channel chan *FilterResponse) {
	isKibana := strings.HasPrefix(kfr.info.Name(), ".kibana")

	if isKibana {
		kfr.report.AddReason("Index '%v' is a kibana index", kfr.info.Name())
	} else {
		kfr.report.AddReason("Index '%v' is  not a kibana index", kfr.info.Name())
	}

	channel <- kfr.response.setPassed(isKibana && kfr.criteria.Include()).setReport(kfr.report)
}
