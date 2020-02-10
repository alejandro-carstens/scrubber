package runners

import (
	"errors"

	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/criterias"
	"github.com/alejandro-carstens/scrubber/actions/filters/runners/reports"
	"github.com/alejandro-carstens/scrubber/actions/infos"
)

type baseRunner struct {
	info       infos.Informable
	report     *reports.Report
	connection *golastic.Connection
}

// BaseInit initializes the base properties for a filter runner
func (br *baseRunner) BaseInit(
	criteria criterias.Criteriable,
	connection *golastic.Connection,
	info ...infos.Informable,
) error {
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

	report.SetCriteria(criteria)

	br.report = report
	br.connection = connection

	return nil
}
