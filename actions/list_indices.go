package actions

import (
	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

type listIndices struct {
	filterAction
	options *options.ListIndicesOptions
}

// ApplyOptions implementation of the Actionable interface
func (li *listIndices) ApplyOptions() Actionable {
	li.options = li.context.Options().(*options.ListIndicesOptions)

	li.indexer.SetOptions(&golastic.IndexOptions{Timeout: li.options.TimeoutInSeconds()})

	return li
}

// Perform implementation of the Actionable interface
func (li *listIndices) Perform() Actionable {
	li.reporter.logger.Noticef("indices:")

	for _, element := range li.list {
		li.reporter.logger.Noticef("%v", element)
	}

	return li
}
