package actions

import (
	"scrubber/actions/contexts"
	"scrubber/logging"

	"github.com/alejandro-carstens/golastic"
)

type action struct {
	retryCount     int
	name           string
	context        contexts.Contextable
	builder        *golastic.ElasticsearchBuilder
	reporter       *reporter
	errorReportMap *errorReportMap
}

func (a *action) Init(context contexts.Contextable, logger *logging.SrvLogger) error {
	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		return err
	}

	a.builder = builder
	a.context = context
	a.name = context.Action()
	a.errorReportMap = newErrorReportMap()
	a.reporter = newReporter(logger)

	return nil
}

func (a *action) HasErrors() bool {
	for _, report := range a.errorReportMap.reports {
		if a.reporter.Logger() != nil {
			for _, err := range report.errs {
				a.reporter.Logger().Errorf("Errors: %v", err)
			}

			continue
		}

		break
	}

	return a.errorReportMap.hasErrors()
}

func (a *action) DisableAction() bool {
	return a.context.Options().GetDisableAction()
}
