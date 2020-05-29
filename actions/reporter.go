package actions

import (
	"scrubber/actions/filters/runners/reports"
	"scrubber/logger"
)

type reporter struct {
	reports []reports.Reportable
	logger  *logger.Logger
}

func (r *reporter) addReports(reports ...reports.Reportable) *reporter {
	if len(r.reports) == 0 {
		r.reports = reports

		return r
	}

	r.reports = append(r.reports, reports...)

	return r
}

func (r *reporter) logFilterResults() error {
	for _, report := range r.reports {
		line, err := report.Line()

		if err != nil {
			return err
		}

		r.logger.Debugf("%v", line)
	}

	return nil
}
