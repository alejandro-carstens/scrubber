package actions

import (
	"errors"
	"github.com/alejandro-carstens/scrubber/actions/options"

	"github.com/alejandro-carstens/golastic"
)

type closeIndices struct {
	filterAction
	options *options.CloseIndicesOptions
}

// ApplyOptions implementation of the Actionable interface
func (ci *closeIndices) ApplyOptions() Actionable {
	ci.options = ci.context.Options().(*options.CloseIndicesOptions)

	ci.indexer.SetOptions(&golastic.IndexOptions{Timeout: ci.options.TimeoutInSeconds()})

	return ci
}

// Perform implementation of the Actionable interface
func (ci *closeIndices) Perform() Actionable {
	ci.exec(func(index string) error {
		response, err := ci.indexer.Close(index)

		if err != nil {
			return err
		}

		if acknowledged, _ := response.S("acknowledged").Data().(bool); !acknowledged {
			return errors.New("close action was not acknowledged")
		}

		return nil
	})

	return ci
}
