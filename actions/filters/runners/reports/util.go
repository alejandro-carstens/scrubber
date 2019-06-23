package reports

// NewAggregateReport returns an instance of *AggregateReport
func NewAggregateReport() *AggregateReport {
	return new(AggregateReport)
}

// NewReport returns an instance of *Report
func NewReport() *Report {
	return new(Report)
}
