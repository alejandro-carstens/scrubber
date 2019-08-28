package runners

import (
	"scrubber/actions/criterias"
	"scrubber/actions/infos"

	"github.com/alejandro-carstens/golastic"
)

type aliasFilterRunner struct {
	baseRunner
}

// Init initializes the filter runner
func (afr *aliasFilterRunner) Init(connection *golastic.Connection, info ...infos.Informable) (Runnerable, error) {
	err := afr.BaseInit(connection, info...)

	return afr, err
}

// RunFilter filters out elements from the actionable list
func (afr *aliasFilterRunner) RunFilter(channel chan *FilterResponse, criteria criterias.Criteriable) {
	if err := afr.validateCriteria(criteria); err != nil {
		channel <- afr.response.setError(err)

		return
	}

	alias := criteria.(*criterias.Alias)

	container, err := afr.connection.Indexer(nil).AliasesCat()

	if err != nil {
		channel <- afr.response.setError(err)

		return
	}

	aliasesResponse, err := container.Children()

	if err != nil {
		channel <- afr.response.setError(err)

		return
	}

	passed := false

	for _, aliasResponse := range aliasesResponse {
		indexName, _ := aliasResponse.S("index").Data().(string)
		aliasName, _ := aliasResponse.S("alias").Data().(string)

		if afr.info.Name() == indexName && afr.inSlice(aliasName, alias.Aliases) {
			afr.report.AddReason("Alias '%v' matched for index '%v'", aliasName, afr.info.Name())

			passed = true
			break
		}
	}

	if !passed {
		afr.report.AddReason("Alias not matched for index '%v'", afr.info.Name())
	}

	channel <- afr.response.setPassed(passed && alias.Include()).setReport(afr.report)
}

func (afr *aliasFilterRunner) inSlice(needle string, haystack []string) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}

	return false
}
