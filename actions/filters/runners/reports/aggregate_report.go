package reports

import "fmt"

// AggregateReport represents an aggregate filter report
type AggregateReport struct {
	baseReport
	Names   []string `json:"names"`
	Results []string `json:"results"`
}

// AddName appends actionable list names to the report
func (ar *AggregateReport) AddName(name string) *AggregateReport {
	if len(ar.Names) == 0 {
		ar.Names = []string{}
	}

	ar.Names = append(ar.Names, name)

	return ar
}

// AddResults appends report results
func (ar *AggregateReport) AddResults(results ...string) *AggregateReport {
	if len(ar.Results) == 0 {
		ar.Results = []string{}
	}

	ar.Results = append(ar.Results, results...)

	return ar
}

// Error adds an error message to the summary
func (ar *AggregateReport) Error(err error) Reportable {
	ar.Summary = append([]string{}, err.Error())

	return ar
}

// Line returns a human readable string for a filter runner's activity
func (ar *AggregateReport) Line() (string, error) {
	criteria, err := ar.toJsonString(ar.Criteria)

	if err != nil {
		return "", err
	}

	summary := "\n"

	for i, reason := range ar.Summary {
		if i+1 == len(ar.Summary) {
			summary = summary + reason
			break
		}

		summary = summary + reason + "\n"
	}

	return fmt.Sprintf(
		"\nType: %v\nNames: %v\nFilter Type: %v\nActionable List: %v\nCriteria: %v\nSummary: %v\n",
		ar.Type,
		ar.Names,
		ar.FilterType,
		ar.Results,
		criteria,
		summary,
	), nil
}
