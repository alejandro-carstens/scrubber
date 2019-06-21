package runners

import (
	"errors"
	"scrubber/actions/criterias"
	"scrubber/actions/filters/runners/reports"
	"scrubber/actions/responses"
)

type baseRunner struct {
	info     responses.Informable
	report   *reports.Report
	response *FilterResponse
}

func (br *baseRunner) BaseInit(info ...responses.Informable) error {
	if len(info) != 1 {
		return errors.New("This is not an aggregate filter runner and as such only accepts one index per run")
	}

	br.info = info[0]

	if len(br.info.Name()) == 0 {
		return errors.New("Could not retrieve name from info")
	}

	report := reports.NewReport().SetName(br.info.Name())

	if br.info.IsSnapshotInfo() {
		report.SetType("snapshot")
	} else {
		report.SetType("index")
	}

	br.report = report
	br.response = new(FilterResponse)

	return nil
}

func (br *baseRunner) validateCriteria(criteria criterias.Criteriable) error {
	br.report.SetCriteria(criteria)

	return criteria.Validate()
}
