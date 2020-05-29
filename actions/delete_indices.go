package actions

import (
	"github.com/alejandro-carstens/golastic"
	"scrubber/actions/options"
)

type deleteIndices struct {
	filterAction
	options *options.DeleteIndicesOptions
}

// ApplyOptions implementation of the Actionable interface
func (di *deleteIndices) ApplyOptions() Actionable {
	di.options = di.context.Options().(*options.DeleteIndicesOptions)

	di.indexer.SetOptions(&golastic.IndexOptions{Timeout: di.options.TimeoutInSeconds()})

	return di
}

// Perform implementation of the Actionable interface
func (di *deleteIndices) Perform() Actionable {
	di.exec(func(index string) error {
		return di.indexer.DeleteIndex(index)
	})

	return di
}
