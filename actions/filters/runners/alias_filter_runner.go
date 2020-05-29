package runners

import (
	"scrubber/actions/criterias"
	"scrubber/actions/infos"

	"github.com/alejandro-carstens/golastic"
)

type aliasFilterRunner struct {
	baseRunner
	criteria *criterias.Alias
}

// Init initializes the filter runner
func (afr *aliasFilterRunner) Init(
	criteria criterias.Criteriable,
	connection *golastic.Connection,
	info ...infos.Informable,
) (Runnerable, error) {
	if err := afr.BaseInit(criteria, connection, info...); err != nil {
		return nil, err
	}

	afr.criteria = criteria.(*criterias.Alias)

	return afr, nil
}

// RunFilter filters out elements from the actionable list
func (afr *aliasFilterRunner) RunFilter(channel chan *FilterResponse) {
	container, err := afr.connection.Indexer(nil).AliasesCat()

	if err != nil {
		channel <- &FilterResponse{Err: err}

		return
	}

	aliasesResponse, err := container.Children()

	if err != nil {
		channel <- &FilterResponse{Err: err}

		return
	}

	passed := false

	for _, aliasResponse := range aliasesResponse {
		indexName, _ := aliasResponse.S("index").Data().(string)
		aliasName, _ := aliasResponse.S("alias").Data().(string)

		if afr.info.Name() == indexName && inStringSlice(aliasName, afr.criteria.Aliases) {
			afr.report.AddReason("Alias '%v' matched for index '%v'", aliasName, afr.info.Name())

			passed = true
			break
		}
	}

	if !passed {
		afr.report.AddReason("Alias not matched for index '%v'", afr.info.Name())
	}

	channel <- &FilterResponse{
		Err:    err,
		Passed: passed && afr.criteria.Include(),
		Report: afr.report,
	}
}
