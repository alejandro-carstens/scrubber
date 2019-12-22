package actions

import (
	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

type listSnapshots struct {
	filterAction
	options *options.ListSnapshotsOptions
}

// ApplyOptions implementation of the Actionable interface
func (ls *listSnapshots) ApplyOptions() Actionable {
	ls.options = ls.context.Options().(*options.ListSnapshotsOptions)

	ls.indexer.SetOptions(&golastic.IndexOptions{Timeout: ls.options.TimeoutInSeconds()})

	return ls
}

// Perform implementation of the Actionable interface
func (ls *listSnapshots) Perform() Actionable {
	ls.reporter.Logger().Noticef("snapshots:")

	for _, element := range ls.list {
		ls.reporter.Logger().Noticef("%v", element)
	}

	return ls
}
