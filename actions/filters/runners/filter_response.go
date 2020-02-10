package runners

import "github.com/alejandro-carstens/scrubber/actions/filters/runners/reports"

// FilterResponse is the response of a filter run
type FilterResponse struct {
	Passed bool               `json:"passed"`
	Err    error              `json:"error"`
	List   []string           `json:"list"`
	Report reports.Reportable `json:"report"`
}
