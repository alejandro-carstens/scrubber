package actions

import (
	"scrubber/actions/contexts"
	"scrubber/logger"

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

// Init initializes an action 
func (a *action) Init(context contexts.Contextable, logger *logger.Logger) error {
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

// HasErrors signals and logs whether the action experienced errors
func (a *action) HasErrors() bool {
	for _, report := range a.errorReportMap.reports {
		if a.reporter.Logger() == nil {
			break
		}

		for _, err := range report.errs {
			a.reporter.Logger().Errorf("Errors: %v", err)
		}
	}

	return a.errorReportMap.hasErrors()
}

// DisableAction indicates whether or not the action should be performed
func (a *action) DisableAction() bool {
	return a.context.Options().GetDisableAction()
}

// List returns the actionable list
func (a *action) List() []string {
	return []string{}
}
