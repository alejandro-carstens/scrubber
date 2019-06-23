package reports

import "fmt"

// Report represents a filter report
type Report struct {
	baseReport
	Name   string `json:"name"`
	Result bool   `json:"result"`
}

// Error adds an error message to the summary
func (r *Report) Error(err error) Reportable {
	if len(r.Summary) == 0 {
		r.Summary = append([]string{}, err.Error())

		return r
	}

	r.Summary = append(r.Summary, err.Error())

	return r
}

// SetName sets the name of the element being filtered
func (r *Report) SetName(name string) *Report {
	r.Name = name

	return r
}

// SetResult sets the result of a given filter
func (r *Report) SetResult(result bool) *Report {
	r.Result = result

	return r
}

// Line returns a human readable string or a filter runner activity
func (r *Report) Line() (string, error) {
	criteria, err := r.toJsonString(r.Criteria)

	if err != nil {
		return "", err
	}

	summary := "\n"

	for i, reason := range r.Summary {
		if i+1 == len(r.Summary) {
			summary = summary + reason

			break
		}

		summary = summary + reason + "\n"
	}

	return fmt.Sprintf(
		"\nType: %v\nName: %v\nFilter Type: %v\nPassed Filter: %t\nCriteria: %v\nSummary: %v\n",
		r.Type,
		r.Name,
		r.FilterType,
		r.Result,
		criteria,
		summary,
	), nil
}
