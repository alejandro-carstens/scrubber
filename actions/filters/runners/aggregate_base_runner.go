package runners

import (
	"errors"
	"scrubber/actions/criterias"
	"scrubber/actions/filters/runners/reports"
	"scrubber/actions/infos"
)

type aggregateBaseRunner struct {
	info     map[string]infos.Informable
	report   *reports.AggregateReport
	response *FilterResponse
}

// BaseInit initializes the base properties for a filter runner
func (abr *aggregateBaseRunner) BaseInit(info ...infos.Informable) error {
	if len(info) == 0 {
		return errors.New("info cannot be empty")
	}

	abr.report = reports.NewAggregateReport()
	abr.response = new(FilterResponse)
	abr.info = map[string]infos.Informable{}

	for _, element := range info {
		if len(element.Name()) == 0 {
			return errors.New("Could not retrieve name")
		}

		abr.info[element.Name()] = element
		abr.report.AddName(element.Name())
	}

	return nil
}

func (abr *aggregateBaseRunner) validateCriteria(criteria criterias.Criteriable) error {
	abr.report.SetCriteria(criteria)

	return criteria.Validate()
}

func (abr *aggregateBaseRunner) excludeIndices(exludeList []string) []string {
	list := []string{}

	elementMap := abr.info

	for _, element := range exludeList {
		delete(elementMap, element)
	}

	for name, _ := range elementMap {
		list = append(list, name)
	}

	return list
}
