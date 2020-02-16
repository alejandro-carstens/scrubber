package filters

import (
	"github.com/alejandro-carstens/scrubber/actions/filters/runners"
	"github.com/alejandro-carstens/scrubber/actions/filters/runners/reports"
	"github.com/alejandro-carstens/scrubber/actions/infos"
)

type FilterRunner struct {
	baseFilterRunner
	info infos.Informable
}

func (fr *FilterRunner) ApplyFilters() (bool, error) {
	channel := make(chan *runners.FilterResponse, len(fr.builder.Criteria()))

	for _, criteria := range fr.builder.Criteria() {
		runner, err := runners.NewRunner(criteria, fr.connection, fr.info)

		if err != nil {
			return false, err
		}

		go runner.RunFilter(channel)
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

	fr.release(channel)

	return true, nil
}
