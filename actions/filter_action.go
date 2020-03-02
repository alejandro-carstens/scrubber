package actions

import (
	"github.com/alejandro-carstens/scrubber/actions/filters"
	"github.com/alejandro-carstens/scrubber/actions/infos"
	"github.com/ivpusic/grpool"
)

type runAction func(element string) error

type filterAction struct {
	action
	list []string
	info map[string]infos.Informable
}

// ApplyFilters runs filters for each element on the actionable list
func (fa *filterAction) ApplyFilters() error {
	if len(fa.context.ActionableList()) > 0 {
		fa.list = append([]string{}, fa.context.ActionableList()...)

		return nil
	} else if fa.context.Options().IsSnapshot() && fa.context.Options().Exists("name") {
		return fa.setSanpshotInfo()
	}

	var actionableList []string
	var err error

	if fa.context.Options().IsSnapshot() {
		actionableList, err = fa.indexer.ListSnapshots(fa.context.Options().String("repository"))
	} else {
		actionableList, err = fa.indexer.ListAllIndices()
	}

	if err != nil {
		return err
	}

	if fa.context.Builder().IsEmpty() {
		fa.list = actionableList

		return nil
	}

	list := []string{}

	for _, element := range actionableList {
		include, err := fa.runFilter(element)

		if err != nil {
			fa.reporter.logFilterResults()

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

	fa.reporter.logFilterResults()

	return err
}

// List returns the actionable list
func (fa *filterAction) List() []string {
	return fa.list
}

func (fa *filterAction) runFilter(element string) (bool, error) {
	var info infos.Informable
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
		fa.info = map[string]infos.Informable{}
	}

	fa.info[element] = info

	runner, err := filters.NewFilterRunner(info, fa.context.Builder(), fa.connection)

	if err != nil {
		return false, err
	}

	passed, err := runner.ApplyFilters()

	fa.reporter.addReports(runner.GetReports()...)

	return passed, err
}

func (fa *filterAction) runAggregateFilters(list []string) ([]string, error) {
	info := []infos.Informable{}

	for _, name := range list {
		info = append(info, fa.info[name])
	}

	runner, err := filters.NewAggregateFilterRunner(info, fa.context.Builder(), fa.connection)

	if err != nil {
		return nil, err
	}

	indicesList, err := runner.ApplyFilters()

	fa.reporter.addReports(runner.GetReports()...)

	return indicesList, err
}

func (fa *filterAction) exec(fn runAction) {
	if fa.context.GetAsync() {
		pool := grpool.NewPool(fa.context.GetNumberOfWorkers(), fa.context.GetQueueLength())

		defer pool.Release()

		pool.WaitCount(len(fa.list))

		for _, element := range fa.list {
			index := element

			pool.JobQueue <- func() {
				defer pool.JobDone()

				if err := fn(index); err != nil {
					fa.errorContainer.push(fa.name, index, err)

					return
				}

				fa.notifiableList = append(fa.notifiableList, index)
			}
		}

		pool.WaitAll()
	} else {
		for _, index := range fa.list {
			if err := fn(index); err != nil {
				fa.errorContainer.push(fa.name, index, err)

				continue
			}

			fa.notifiableList = append(fa.notifiableList, index)
		}
	}

	if len(fa.errorContainer.list()) > 0 && fa.retryCount < fa.context.GetRetryCount() {
		fa.retryCount = fa.retryCount + 1
		fa.list = fa.errorContainer.list()

		fa.exec(fn)
	}
}

func (fa *filterAction) fetchSnapshot(snapshot string) (infos.Informable, error) {
	response, err := fa.indexer.GetSnapshots(fa.context.Options().String("repository"), snapshot)

	if err != nil {
		return nil, err
	}

	children, err := response.S("snapshots").Children()

	if err != nil {
		return nil, err
	}

	return new(infos.SnapshotInfo).Marshal(children[0])
}

func (fa *filterAction) fetchIndexCat(index string) (infos.Informable, error) {
	response, err := fa.indexer.IndexCat(index)

	if err != nil {
		return nil, err
	}

	children, err := response.Children()

	if err != nil {
		return nil, err
	}

	return new(infos.IndexInfo).Marshal(children[0])
}

func (fa *filterAction) setSanpshotInfo() error {
	name := fa.context.Options().String("name")

	fa.list = append([]string{}, name)

	info, err := fa.fetchSnapshot(name)

	if err != nil {
		return err
	}

	fa.info = map[string]infos.Informable{
		name: info,
	}

	return nil
}
