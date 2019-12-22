package actions

import (
	"errors"

	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

type indexSettings struct {
	filterAction
	options *options.IndexSettingsOptions
}

// ApplyOptions implementation of the Actionable interface
func (is *indexSettings) ApplyOptions() Actionable {
	is.options = is.context.Options().(*options.IndexSettingsOptions)

	is.indexer.SetOptions(&golastic.IndexOptions{Timeout: is.options.TimeoutInSeconds()})

	return is
}

// Perform implementation of the Actionable interface
func (is *indexSettings) Perform() Actionable {
	is.exec(func(index string) error {
		settings, err := mapToString(is.options.IndexSettings)

		if err != nil {
			return err
		}

		response, err := is.indexer.PutSettings(settings, index)

		if err != nil {
			return err
		}

		if acknowledged, _ := response.S("acknowledged").Data().(bool); !acknowledged {
			return errors.New("index settings action was not acknowledged")
		}

		return nil
	})

	return is
}
