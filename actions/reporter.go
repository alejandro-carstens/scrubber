package actions

import (
	"fmt"
	"scrubber/actions/filters/runners/reports"
	"scrubber/logger"
)

type reporter struct {
	reports []reports.Reportable
	logger  *logger.Logger
}

func (r *reporter) AddReports(reports ...reports.Reportable) *reporter {
	if len(r.reports) == 0 {
		r.reports = reports

		return r
	}

	r.reports = append(r.reports, reports...)

	return r
}

func (r *reporter) Logger() *logger.Logger {
	return r.logger
}

func (r *reporter) LogFilterResults() error {
	for _, report := range r.reports {
		line, err := report.Line()

		if err != nil {
			return err
		}

		if r.logger != nil {
			r.logger.Debugf("%v", line)
		} else {
			fmt.Println(line)
		}
	}

	return nil
}
