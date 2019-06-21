package runners

import (
	"scrubber/actions/criterias"
	"scrubber/actions/responses"

	"github.com/alejandro-carstens/golastic"
)

type allocatedFilterRunner struct {
	baseRunner
	builder golastic.Queryable
}

func (afr *allocatedFilterRunner) Init(info ...responses.Informable) (Runnerable, error) {
	if err := afr.BaseInit(info...); err != nil {
		return nil, err
	}

	builder, err := golastic.NewBuilder(golastic.NewGolasticModel(), nil)

	if err != nil {
		return nil, err
	}

	afr.builder = builder

	return afr, nil
}

func (afr *allocatedFilterRunner) RunFilter(channel chan *FilterResponse, criteria criterias.Criteriable) {
	if err := afr.validateCriteria(criteria); err != nil {
		channel <- afr.response.setError(err)
		return
	}

	allocated := criteria.(*criterias.Allocated)

	settingsResponse, err := afr.builder.Settings(afr.info.Name())

	if err != nil {
		channel <- afr.response.setError(err)
		return
	}

	passed := false
	container := settingsResponse[afr.info.Name()]

	value, valid := container.S("index", "routing", "allocation", allocated.AllocationType, allocated.Key).Data().(string)

	if valid {
		passed = value == allocated.Value

		if passed {
			afr.report.AddReason(
				"Allocation value '%v' matched for allocation type '%v' and key '%v'",
				value,
				allocated.AllocationType,
				allocated.Key,
			)
		}
	}

	if !passed {
		afr.report.AddReason(
			"Allocation value '%v' not matched for allocation type '%v' and key '%v'",
			value,
			allocated.AllocationType,
			allocated.Key,
		)
	}

	channel <- afr.response.setPassed(passed && allocated.Include()).setReport(afr.report)
}
