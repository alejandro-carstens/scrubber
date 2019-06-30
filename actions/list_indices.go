package actions

import (
	"scrubber/actions/options"

	"github.com/alejandro-carstens/golastic"
)

type listIndices struct {
	filterAction
	options *options.ListIndicesOptions
}

func (li *listIndices) ApplyOptions() Actionable {
	li.options = li.context.Options().(*options.ListIndicesOptions)

	li.builder.SetOptions(&golastic.IndexOptions{Timeout: li.options.TimeoutInSeconds()})

	return li
}

func (li *listIndices) Perform() Actionable {
	li.reporter.Logger().Noticef("indices:")

	for _, element := range li.list {
		li.reporter.Logger().Noticef("%v", element)
	}

	return li
}
