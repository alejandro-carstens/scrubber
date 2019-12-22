package runners

import "github.com/alejandro-carstens/scrubber/actions/filters/runners/reports"

// FilterResponse is the response of a filter run
type FilterResponse struct {
	Passed bool               `json:"passed"`
	Err    error              `json:"error"`
	List   []string           `json:"list"`
	Report reports.Reportable `json:"report"`
}

func (fr *FilterResponse) setError(err error) *FilterResponse {
	fr.Err = err

	return fr
}

func (fr *FilterResponse) setPassed(passed bool) *FilterResponse {
	fr.Passed = passed

	return fr
}

func (fr *FilterResponse) setReport(report reports.Reportable) *FilterResponse {
	fr.Report = report

	return fr
}

func (fr *FilterResponse) setList(list []string) *FilterResponse {
	fr.List = list

	return fr
}
