package runners

import (
	"regexp"
	"github.com/alejandro-carstens/scrubber/actions/criterias"
	"github.com/alejandro-carstens/scrubber/actions/infos"
	"sort"

	"github.com/alejandro-carstens/golastic"
)

type countFilterRunner struct {
	useAgeFilterRunner
	criteria *criterias.Count
}

// Init initializes the filter runner
func (cfr *countFilterRunner) Init(criteria criterias.Criteriable, connection *golastic.Connection, info ...infos.Informable) (Runnerable, error) {
	if err := cfr.BaseInit(criteria, connection, info...); err != nil {
		return nil, err
	}

	cfr.criteria = criteria.(*criterias.Count)

	return cfr, nil
}

// RunFilter filters out elements from the actionable list
func (cfr *countFilterRunner) RunFilter(channel chan *FilterResponse) {
	var err error
	var sortedList []string

	if cfr.criteria.UseAge {
		sortedList, err = cfr.runAgeSorters(cfr.criteria)
	} else if len(cfr.criteria.Pattern) > 0 {
		sortedList, err = cfr.sortByPattern()
	} else {
		sortedList = cfr.runDefaultSorter(cfr.criteria)
	}

	if err == nil && len(sortedList) >= cfr.criteria.Count {
		sortedList = sortedList[:cfr.criteria.Count]
	}

	if !cfr.criteria.Include() {
		cfr.report.AddReason("Excluding indices: %v from list", sortedList)

		sortedList = cfr.excludeIndices(sortedList)
	}

	cfr.report.AddResults(sortedList...)

	channel <- cfr.response.setError(err).setPassed(true).setReport(cfr.report).setList(sortedList)
}

func (cfr *countFilterRunner) sortByPattern() ([]string, error) {
	list := []string{}

	for name, _ := range cfr.info {
		match, err := regexp.MatchString(cfr.criteria.Pattern, name)

		if err != nil {
			return nil, err
		}

		if match {
			list = append(list, name)
		}
	}

	if cfr.criteria.Reverse {
		sort.Sort(sort.Reverse(sort.StringSlice(list)))
	} else {
		sort.Sort(sort.StringSlice(list))
	}

	cfr.report.AddReason("Sorted List: %v", list)

	return list, nil
}
