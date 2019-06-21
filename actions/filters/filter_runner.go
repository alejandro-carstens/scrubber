package filters

import (
	"scrubber/actions/filters/runners"
	"scrubber/actions/filters/runners/reports"
	"scrubber/actions/responses"
)

type FilterRunner struct {
	baseFilterRunner
	info responses.Informable
}

func (fr *FilterRunner) ApplyFilters() (bool, error) {
	channel := make(chan *runners.FilterResponse, len(fr.builder.Criteria()))

	for _, criteria := range fr.builder.Criteria() {
		runner, err := runners.NewRunner(criteria.Name(), fr.info)

		if err != nil {
			return false, err
		}

		go runner.RunFilter(channel, criteria)
	}

	for range fr.builder.Criteria() {
		filterResponse := <-channel

		if filterResponse.Err != nil {
			fr.AddReport(filterResponse.Report.Error(filterResponse.Err))

			return false, filterResponse.Err
		}

		if !filterResponse.Passed {
			fr.AddReport(filterResponse.Report)

			return false, nil
		}

		fr.AddReport(filterResponse.Report.(*reports.Report).SetResult(true))
	}

	return true, nil
}
