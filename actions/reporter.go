package actions

import (
	"fmt"

	"github.com/alejandro-carstens/scrubber/actions/filters/runners/reports"
	"github.com/alejandro-carstens/scrubber/logger"
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

		if r.logger != nil {
			r.logger.Debugf("%v", line)
		} else {
			fmt.Println(line)
		}
	}

	return nil
}
