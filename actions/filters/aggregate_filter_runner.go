package filters

import (
	"scrubber/actions/filters/runners"
	"scrubber/actions/infos"
)

type AggregateFilterRunner struct {
	baseFilterRunner
	info     []infos.Informable
	countMap map[string]int
}

func (afr *AggregateFilterRunner) ApplyFilters() ([]string, error) {
	channel := make(chan *runners.FilterResponse, len(afr.builder.AggregateCriteria()))

	for _, criteria := range afr.builder.AggregateCriteria() {
		runner, err := runners.NewRunner(criteria.Name(), afr.connection, afr.info...)

		if err != nil {
			return nil, err
		}

		go runner.RunFilter(channel, criteria)
	}

	for i := 0; i < len(afr.builder.AggregateCriteria()); i++ {
		filterResponse := <-channel

		if filterResponse.Err != nil || !filterResponse.Passed {
			return nil, filterResponse.Err
		}

		afr.AddReport(filterResponse.Report)
		afr.addIndicesToCountMap(filterResponse.List...)
	}

	return afr.getIndicesList(len(afr.builder.AggregateCriteria())), nil
}

func (afr *AggregateFilterRunner) addIndicesToCountMap(list ...string) {
	for _, element := range list {
		if _, valid := afr.countMap[element]; valid {
			afr.countMap[element] = afr.countMap[element] + 1

			continue
		}

		afr.countMap[element] = 1
	}
}

func (afr *AggregateFilterRunner) getIndicesList(expectedCount int) []string {
	list := []string{}

	for name, count := range afr.countMap {
		if count == expectedCount {
			list = append(list, name)
		}
	}

	return list
}
