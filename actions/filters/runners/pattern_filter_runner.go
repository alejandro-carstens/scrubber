package runners

import (
	"regexp"
	"strings"

	"scrubber/actions/criterias"
	"scrubber/actions/infos"

	"github.com/alejandro-carstens/golastic"
)

type patternFilterRunner struct {
	baseRunner
	criteria *criterias.Pattern
}

// Init initializes the filter runner
func (pfr *patternFilterRunner) Init(
	criteria criterias.Criteriable,
	connection *golastic.Connection,
	info ...infos.Informable,
) (Runnerable, error) {
	if err := pfr.BaseInit(criteria, connection, info...); err != nil {
		return nil, err
	}

	pfr.criteria = criteria.(*criterias.Pattern)

	return pfr, nil
}

// RunFilter filters out elements from the actionable list
func (pfr *patternFilterRunner) RunFilter(channel chan *FilterResponse) {
	var passed bool
	var err error

	switch pfr.criteria.Kind {
	case "timestring":
		passed, err = pfr.processTimestring()
		break
	case "prefix":
		passed, err = pfr.processPrefix()
		break
	case "suffix":
		passed, err = pfr.processSuffix()
		break
	case "regex":
		passed, err = pfr.processRegex()
		break
	}

	channel <- &FilterResponse{
		Err:    err,
		Passed: passed && pfr.criteria.Include(),
		Report: pfr.report,
	}
}

func (pfr *patternFilterRunner) processTimestring() (bool, error) {
	var regPattern string

	switch pfr.criteria.Value {
	case "Y.m.d":
		regPattern = `\d{4}.\d{2}.\d{2}`
		break
	case "m.d.Y":
		regPattern = `\d{2}.\d{2}.\d{4}`
		break
	case "Y.m":
		regPattern = `\d{4}.\d{2}`
		break
	case "Y-m-d":
		regPattern = `\d{4}-\d{2}-\d{2}`
		break
	case "Y-m-d H:M":
		regPattern = `\d{4}-\d{2}-\d{2} \d{2}:\d{2}`
		break
	case "Y-m-d H:M:S":
		regPattern = `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`
		break
	case "m-d-Y":
		regPattern = `\d{2}-\d{2}-\d{4}`
		break
	case "Y-m":
		regPattern = `\d{4}-\d{2}`
		break
	}

	reg, err := regexp.Compile(regPattern)

	if err != nil {
		return false, nil
	}

	pfr.report.AddReason("Processing timestring pattern '%v'", pfr.criteria.Value)

	if passed := len(reg.FindString(pfr.info.Name())) > 0; passed {
		pfr.report.AddReason("Index '%v' matched pattern '%v'", pfr.info.Name(), pfr.criteria.Value)

		return true, nil
	}

	pfr.report.AddReason("Index '%v' did not matched pattern '%v'", pfr.info.Name(), pfr.criteria.Value)

	return false, nil
}

func (pfr *patternFilterRunner) processSuffix() (bool, error) {
	pfr.report.AddReason("Processing by suffix pattern '%v'", pfr.criteria.Value)

	if passed := strings.HasSuffix(pfr.info.Name(), pfr.criteria.Value); passed {
		pfr.report.AddReason("Index '%v' matched pattern '%v'", pfr.info.Name(), pfr.criteria.Value)

		return passed, nil
	}

	pfr.report.AddReason("Index '%v' did not matched pattern '%v'", pfr.info.Name(), pfr.criteria.Value)

	return false, nil
}

func (pfr *patternFilterRunner) processPrefix() (bool, error) {
	pfr.report.AddReason("Processing by prefix pattern '%v'", pfr.criteria.Value)

	if passed := strings.HasPrefix(pfr.info.Name(), pfr.criteria.Value); passed {
		pfr.report.AddReason("Index '%v' matched pattern '%v'", pfr.info.Name(), pfr.criteria.Value)

		return passed, nil
	}

	pfr.report.AddReason("Index '%v' did not matched pattern '%v'", pfr.info.Name(), pfr.criteria.Value)

	return false, nil
}

func (pfr *patternFilterRunner) processRegex() (bool, error) {
	passed, err := regexp.MatchString(pfr.criteria.Value, pfr.info.Name())

	if err != nil {
		return false, err
	}

	pfr.report.AddReason("Processing by regex pattern '%v'", pfr.criteria.Value)

	if passed {
		pfr.report.AddReason("Index '%v' matched pattern '%v'", pfr.info.Name(), pfr.criteria.Value)

		return passed, nil
	}

	pfr.report.AddReason("Index '%v' did not matched pattern '%v'", pfr.info.Name(), pfr.criteria.Value)

	return false, nil
}
