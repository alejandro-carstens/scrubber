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
func (afr *aliasFilterRunner) Init(builder *golastic.ElasticsearchBuilder, info ...infos.Informable) (Runnerable, error) {
	if err := afr.BaseInit(builder, info...); err != nil {
		return nil, err
	}

	return afr, nil
}

// RunFilter filters out elements from the actionable list
func (afr *aliasFilterRunner) RunFilter(channel chan *FilterResponse, criteria criterias.Criteriable) {
	if err := afr.validateCriteria(criteria); err != nil {
		channel <- afr.response.setError(err)
		return
	}

	alias := criteria.(*criterias.Alias)

	aliasesResponse, err := afr.builder.AliasesCat()

	if err != nil {
		channel <- afr.response.setError(err)
		return
	}

	passed := false

	for _, aliasResponse := range aliasesResponse {
		if afr.info.Name() == aliasResponse.Index && afr.inSlice(aliasResponse.Alias, alias.Aliases) {
			afr.report.AddReason(
				"Alias '%v' matched for index '%v'",
				aliasResponse.Alias,
				afr.info.Name(),
			)

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
