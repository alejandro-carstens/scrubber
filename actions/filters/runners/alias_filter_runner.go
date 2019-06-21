package runners

import (
	"scrubber/actions/criterias"
	"scrubber/actions/responses"

	"github.com/alejandro-carstens/golastic"
)

type aliasFilterRunner struct {
	baseRunner
	builder golastic.Queryable
}

func (afr *aliasFilterRunner) Init(info ...responses.Informable) (Runnerable, error) {
	if err := afr.BaseInit(info...); err != nil {
		return nil, err
	}

	model := golastic.NewGolasticModel()
	model.SetIndex(afr.info.Name())

	builder, err := golastic.NewBuilder(model, nil)

	if err != nil {
		return nil, err
	}

	afr.builder = builder

	return afr, nil
}

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
