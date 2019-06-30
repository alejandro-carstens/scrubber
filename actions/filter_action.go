package actions

import (
	"scrubber/actions/filters"
	"scrubber/actions/responses"

	"github.com/ivpusic/grpool"
)

type runAction func(element string) error

type filterAction struct {
	action
	info map[string]responses.Informable
	list []string
}

func (fa *filterAction) ApplyFilters() error {
	if len(fa.context.ActionableList()) > 0 {
		fa.list = append([]string{}, fa.context.ActionableList()...)

		return nil
	} else if fa.context.Options().IsSnapshot() && fa.context.Options().Exists("name") {
		fa.list = append([]string{}, fa.context.Options().String("name"))

		return nil
	}

	var actionableList []string
	var err error

	if fa.context.Options().IsSnapshot() {
		actionableList, err = fa.builder.ListSnapshots(fa.context.Options().String("repository"))
	} else {
		actionableList, err = fa.builder.ListAllIndices()
	}

	if err != nil {
		return err
	}

	list := []string{}

	for _, element := range actionableList {
		include, err := fa.runFilter(element)

		if err != nil {
			fa.reporter.LogFilterResults()

			return err
		}

		if include {
			list = append(list, element)
		}
	}

	if len(fa.context.Builder().AggregateCriteria()) > 0 {
		list, err = fa.runAggregateFilters(list)
	}

	fa.list = list

	fa.reporter.LogFilterResults()

	return err
}

func (fa *filterAction) List() []string {
	return fa.list
}

func (fa *filterAction) runFilter(element string) (bool, error) {
	var info responses.Informable
	var err error

	if fa.context.Options().IsSnapshot() {
		info, err = fa.fetchSnapshot(element)
	} else {
		info, err = fa.fetchIndexCat(element)
	}

	if err != nil {
		return false, err
	}

	if len(fa.info) == 0 {
		fa.info = map[string]responses.Informable{}
	}

	fa.info[element] = info

	runner, err := filters.NewFilterRunner(info, fa.context.Builder(), fa.builder, fa.context.Options().IsSnapshot())

	if err != nil {
		return false, err
	}

	passed, err := runner.ApplyFilters()

	fa.reporter.AddReports(runner.GetReports()...)

	return passed, err
}

func (fa *filterAction) runAggregateFilters(list []string) ([]string, error) {
	info := []responses.Informable{}

	for _, name := range list {
		info = append(info, fa.info[name])
	}

	runner, err := filters.NewAggregateFilterRunner(info, fa.context.Builder(), fa.builder)

	if err != nil {
		return nil, err
	}

	indicesList, err := runner.ApplyFilters()

	fa.reporter.AddReports(runner.GetReports()...)

	return indicesList, err
}

func (fa *filterAction) exec(fn runAction) {
	if fa.context.GetAsync() {
		pool := grpool.NewPool(fa.context.GetNumberOfWorkers(), fa.context.GetQueueLength())

		defer pool.Release()

		pool.WaitCount(len(fa.list))

		for _, element := range fa.list {
			param := element

			pool.JobQueue <- func() {
				defer pool.JobDone()

				if err := fn(param); err != nil {
					fa.errorReportMap.push(fa.name, param, err)
				}
			}
		}

		pool.WaitAll()
	} else {
		for _, element := range fa.list {
			if err := fn(element); err != nil {
				fa.errorReportMap.push(fa.name, element, err)
			}
		}
	}

	if len(fa.errorReportMap.list()) > 0 && fa.retryCount < fa.context.GetRetryCount() {
		fa.retryCount = fa.retryCount + 1
		fa.list = fa.errorReportMap.list()

		fa.exec(fn)
	}
}

func (fa *filterAction) fetchSnapshot(snapshot string) (responses.Informable, error) {
	response, err := fa.builder.GetSnapshots(fa.context.Options().String("repository"), snapshot)

	if err != nil {
		return nil, err
	}

	children, err := response.S("snapshots").Children()

	if err != nil {
		return nil, err
	}

	return new(responses.SnapshotInfo).Marshal(children[0])
}

func (fa *filterAction) fetchIndexCat(index string) (responses.Informable, error) {
	response, err := fa.builder.IndexCat(index)

	if err != nil {
		return nil, err
	}

	children, err := response.Children()

	if err != nil {
		return nil, err
	}

	return new(responses.IndexInfo).Marshal(children[0])
}
