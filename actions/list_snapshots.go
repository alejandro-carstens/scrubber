package actions

import (
	"scrubber/actions/options"

	"github.com/alejandro-carstens/golastic"
)

type listSnapshots struct {
	filterAction
	options *options.ListSnapshotsOptions
}

func (ls *listSnapshots) ApplyOptions() Actionable {
	ls.options = ls.context.Options().(*options.ListSnapshotsOptions)

	ls.indexer.SetOptions(&golastic.IndexOptions{Timeout: ls.options.TimeoutInSeconds()})

	return ls
}

func (ls *listSnapshots) Perform() Actionable {
	ls.reporter.Logger().Noticef("snapshots:")

	for _, element := range ls.list {
		ls.reporter.Logger().Noticef("%v", element)
	}

	return ls
}
