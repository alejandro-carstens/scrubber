package actions

type errorReport struct {
	errs       []error
	actionType string
	name       string
	retries    int
	passed     bool
	action     string
}

func (er *errorReport) push(err error) *errorReport {
	er.errs = append(er.errs, err)
	er.retries = er.retries + 1

	return er
}

func (er *errorReport) pass() *errorReport {
	er.passed = true

	return er
}
