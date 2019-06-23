package runners

import (
	"regexp"
	"scrubber/actions/criterias"
	"scrubber/actions/responses"
	"strings"
)

type patternFilterRunner struct {
	baseRunner
}

// Init initializes the filter runner
func (pfr *patternFilterRunner) Init(info ...responses.Informable) (Runnerable, error) {
	err := pfr.BaseInit(info...)

	return pfr, err
}

// RunFilter filters out elements from the actionable list
func (pfr *patternFilterRunner) RunFilter(channel chan *FilterResponse, criteria criterias.Criteriable) {
	if err := pfr.validateCriteria(criteria); err != nil {
		channel <- pfr.response.setError(err)
		return
	}

	pattern := criteria.(*criterias.Pattern)

	var passed bool
	var err error

	switch pattern.Kind {
	case "timestring":
		passed, err = pfr.processTimestring(pattern.Value)
		break
	case "prefix":
		passed, err = pfr.processPrefix(pattern.Value)
		break
	case "suffix":
		passed, err = pfr.processSuffix(pattern.Value)
		break
	case "regex":
		passed, err = pfr.processRegex(pattern.Value)
		break
	}

	channel <- pfr.response.
		setError(err).
		setReport(pfr.report).
		setPassed(passed && pattern.Include())
}

func (pfr *patternFilterRunner) processTimestring(value string) (bool, error) {
	var regPattern string

	switch value {
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

	pfr.report.AddReason("Processing timestring pattern '%v'", value)

	if passed := len(reg.FindString(pfr.info.Name())) > 0; passed {
		pfr.report.AddReason("Index '%v' matched pattern '%v'", pfr.info.Name(), value)

		return true, nil
	}

	pfr.report.AddReason("Index '%v' did not matched pattern '%v'", pfr.info.Name(), value)

	return false, nil
}

func (pfr *patternFilterRunner) processSuffix(value string) (bool, error) {
	pfr.report.AddReason("Processing by suffix pattern '%v'", value)

	if passed := strings.HasSuffix(pfr.info.Name(), value); passed {
		pfr.report.AddReason("Index '%v' matched pattern '%v'", pfr.info.Name(), value)

		return passed, nil
	}

	pfr.report.AddReason("Index '%v' did not matched pattern '%v'", pfr.info.Name(), value)

	return false, nil
}

func (pfr *patternFilterRunner) processPrefix(value string) (bool, error) {
	pfr.report.AddReason("Processing by prefix pattern '%v'", value)

	if passed := strings.HasPrefix(pfr.info.Name(), value); passed {
		pfr.report.AddReason("Index '%v' matched pattern '%v'", pfr.info.Name(), value)

		return passed, nil
	}

	pfr.report.AddReason("Index '%v' did not matched pattern '%v'", pfr.info.Name(), value)

	return false, nil
}

func (pfr *patternFilterRunner) processRegex(value string) (bool, error) {
	passed, err := regexp.MatchString(value, pfr.info.Name())

	if err != nil {
		return false, err
	}

	pfr.report.AddReason("Processing by regex pattern '%v'", value)

	if passed {
		pfr.report.AddReason("Index '%v' matched pattern '%v'", pfr.info.Name(), value)

		return passed, nil
	}

	pfr.report.AddReason("Index '%v' did not matched pattern '%v'", pfr.info.Name(), value)

	return false, nil
}
