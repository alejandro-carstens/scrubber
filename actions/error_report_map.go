package actions

import "fmt"

type errorReportMap struct {
	reports    map[string]*errorReport
	errorCount int
}

func (erm *errorReportMap) push(action string, name string, err error) *errorReportMap {
	errorReport, valid := erm.reports[erm.key(action, name)]

	if valid {
		erm.reports[erm.key(action, name)] = errorReport.push(err)

		return erm
	}

	erm.reports[erm.key(action, name)] = newErrorReport(action, name, err)

	erm.errorCount = erm.errorCount + 1

	return erm
}

func (erm *errorReportMap) clear(action, name string) {
	errorReport, valid := erm.reports[erm.key(action, name)]

	if valid && erm.errorCount > 0 {
		erm.errorCount = erm.errorCount - 1
	}

	erm.reports[erm.key(action, name)] = errorReport.pass()
}

func (erm *errorReportMap) hasErrors() bool {
	return erm.errorCount > 0
}

func (erm *errorReportMap) list() []string {
	list := []string{}

	for _, errorReport := range erm.reports {
		if !errorReport.passed {
			list = append(list, errorReport.name)
		}
	}

	return list
}

func (erm *errorReportMap) key(action, name string) string {
	return fmt.Sprintf("%v-%v", action, name)
}
