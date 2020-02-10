package runners

import (
	"errors"

	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/criterias"
	"github.com/alejandro-carstens/scrubber/actions/filters/runners/reports"
	"github.com/alejandro-carstens/scrubber/actions/infos"
)

type aggregateBaseRunner struct {
	info       map[string]infos.Informable
	report     *reports.AggregateReport
	response   *FilterResponse
	connection *golastic.Connection
}

// BaseInit initializes the base properties for a filter runner
func (abr *aggregateBaseRunner) BaseInit(
	criteria criterias.Criteriable,
	connection *golastic.Connection,
	info ...infos.Informable,
) error {
	if len(info) == 0 {
		return errors.New("info cannot be empty")
	}

	abr.report = reports.NewAggregateReport()
	abr.response = new(FilterResponse)
	abr.info = map[string]infos.Informable{}
	abr.connection = connection

	for i, element := range info {
		if len(element.Name()) == 0 {
			return errors.New("Could not retrieve name")
		}

		if i == 0 {
			if element.IsSnapshotInfo() {
				abr.report.SetType("snapshot")
			} else {
				abr.report.SetType("index")
			}
		}

		abr.info[element.Name()] = element
		abr.report.AddName(element.Name())
	}

	abr.report.SetCriteria(criteria)

	return nil
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
