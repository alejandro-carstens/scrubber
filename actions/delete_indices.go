package actions

import (
	"scrubber/actions/options"

	"github.com/alejandro-carstens/golastic"
)

type deleteIndices struct {
	filterAction
	options *options.DeleteIndicesOptions
}

func (di *deleteIndices) ApplyOptions() Actionable {
	di.options = di.context.Options().(*options.DeleteIndicesOptions)

	di.builder.SetOptions(&golastic.IndexOptions{Timeout: di.options.TimeoutInSeconds()})

	return di
}

func (di *deleteIndices) Perform() Actionable {
	di.exec(func(index string) error {
		return di.builder.DeleteIndex(index)
	})

	return di
}
