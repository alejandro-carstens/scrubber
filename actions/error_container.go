package actions

import "fmt"

type errorReport struct {
	errors     []error
	actionType string
	name       string
	retries    int
	passed     bool
	action     string
}

func (er *errorReport) push(err error) *errorReport {
	er.errors = append(er.errors, err)
	er.retries = er.retries + 1

	return er
}

func (er *errorReport) pass() *errorReport {
	er.passed = true

	return er
}

type errorContainer struct {
	reports    map[string]*errorReport
	errorCount int
}

func (ec *errorContainer) push(action string, name string, err error) *errorContainer {
	errorReport, valid := ec.reports[ec.key(action, name)]

	if valid {
		ec.reports[ec.key(action, name)] = errorReport.push(err)

		return ec
	}

	ec.reports[ec.key(action, name)] = newErrorReport(action, name, err)

	ec.errorCount = ec.errorCount + 1

	return ec
}

func (ec *errorContainer) clear(action, name string) {
	errorReport, valid := ec.reports[ec.key(action, name)]

	if valid && ec.errorCount > 0 {
		ec.errorCount = ec.errorCount - 1
	}

	ec.reports[ec.key(action, name)] = errorReport.pass()
}

func (ec *errorContainer) hasErrors() bool {
	return ec.errorCount > 0
}

func (ec *errorContainer) list() []string {
	list := []string{}

	for _, errorReport := range ec.reports {
		if !errorReport.passed {
			list = append(list, errorReport.name)
		}
	}

	return list
}

func (ec *errorContainer) key(action, name string) string {
	return fmt.Sprintf("%v-%v", action, name)
}
