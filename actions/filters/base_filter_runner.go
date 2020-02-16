package filters

import (
	"errors"

	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/criterias"
	"github.com/alejandro-carstens/scrubber/actions/filters/runners/reports"
)

type baseFilterRunner struct {
	builder    *criterias.Builder
	connection *golastic.Connection
	reports    []reports.Reportable
}

func (bfr *baseFilterRunner) Init(builder *criterias.Builder, connection *golastic.Connection) error {
	if builder == nil {
		return errors.New("Builder can't be nil")
	}

	bfr.builder = builder
	bfr.connection = connection
	bfr.reports = []reports.Reportable{}

	return nil
}

func (bfr *baseFilterRunner) AddReport(report reports.Reportable) {
	bfr.reports = append(bfr.reports, report)
}

func (bfr *baseFilterRunner) GetReports() []reports.Reportable {
	return bfr.reports
}
