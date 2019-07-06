package runners

import (
	"regexp"
	"scrubber/actions/criterias"
	"scrubber/actions/infos"
	"sort"
)

type countFilterRunner struct {
	useAgeFilterRunner
}

// Init initializes the filter runner
func (cfr *countFilterRunner) Init(info ...infos.Informable) (Runnerable, error) {
	if err := cfr.BaseInit(info...); err != nil {
		return nil, err
	}

	return cfr, nil
}

// RunFilter filters out elements from the actionable list
func (cfr *countFilterRunner) RunFilter(channel chan *FilterResponse, criteria criterias.Criteriable) {
	if err := cfr.validateCriteria(criteria); err != nil {
		channel <- cfr.response.setError(err)

		return
	}

	var err error
	var sortedList []string

	count := criteria.(*criterias.Count)

	if count.UseAge {
		sortedList, err = cfr.runAgeSorters(count)
	} else if len(count.Pattern) > 0 {
		sortedList, err = cfr.sortByPattern(count)
	} else {
		sortedList = cfr.runDefaultSorter(count)
	}

	if err == nil && len(sortedList) >= count.Count {
		sortedList = sortedList[:count.Count]
	}

	if !count.Include() {
		cfr.report.AddReason("Excluding indices: %v from list", sortedList)

		sortedList = cfr.excludeIndices(sortedList)
	}

	cfr.report.AddResults(sortedList...)

	channel <- cfr.response.
		setError(err).
		setPassed(true).
		setReport(cfr.report).
		setList(sortedList)
}

func (cfr *countFilterRunner) sortByPattern(count *criterias.Count) ([]string, error) {
	list := []string{}

	for name, _ := range cfr.info {
		match, err := regexp.MatchString(count.Pattern, name)

		if err != nil {
			return nil, err
		}

		if match {
			list = append(list, name)
		}
	}

	if count.Reverse {
		sort.Sort(sort.Reverse(sort.StringSlice(list)))
	} else {
		sort.Sort(sort.StringSlice(list))
	}

	cfr.report.AddReason("Sorted List: %v", list)

	return list, nil
}
