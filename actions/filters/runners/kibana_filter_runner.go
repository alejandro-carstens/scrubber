package runners

import (
	"strings"

	"scrubber/actions/criterias"
	"scrubber/actions/infos"

	"github.com/alejandro-carstens/golastic"
)

type kibanaFilterRunner struct {
	baseRunner
	criteria *criterias.Kibana
}

// Init initializes the filter runner
func (kfr *kibanaFilterRunner) Init(
	criteria criterias.Criteriable,
	connection *golastic.Connection,
	info ...infos.Informable,
) (Runnerable, error) {
	if err := kfr.BaseInit(criteria, connection, info...); err != nil {
		return nil, err
	}

	kfr.criteria = criteria.(*criterias.Kibana)

	return kfr, nil
}

// RunFilter filters out elements from the actionable list
func (kfr *kibanaFilterRunner) RunFilter(channel chan *FilterResponse) {
	passed := strings.HasPrefix(kfr.info.Name(), ".kibana")

	if passed {
		kfr.report.AddReason("Index '%v' is a kibana index", kfr.info.Name())
	} else {
		kfr.report.AddReason("Index '%v' is  not a kibana index", kfr.info.Name())
	}

	channel <- &FilterResponse{
		Passed: passed && kfr.criteria.Include(),
		Report: kfr.report,
	}
}
