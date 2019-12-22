package actions

import (
	"errors"

	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

type openIndices struct {
	filterAction
	options *options.OpenIndicesOptions
}

// ApplyOptions implementation of the Actionable interface
func (oi *openIndices) ApplyOptions() Actionable {
	oi.options = oi.context.Options().(*options.OpenIndicesOptions)

	oi.indexer.SetOptions(&golastic.IndexOptions{Timeout: oi.options.TimeoutInSeconds()})

	return oi
}

// Perform implementation of the Actionable interface
func (oi *openIndices) Perform() Actionable {
	oi.exec(func(index string) error {
		response, err := oi.indexer.Open(index)

		if err != nil {
			return err
		}

		if acknowledged, _ := response.S("acknowledged").Data().(bool); !acknowledged {
			return errors.New("close action was not acknowledged")
		}

		return nil
	})

	return oi
}
