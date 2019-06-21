package actions

import (
	"errors"
	"scrubber/actions/options"

	"github.com/alejandro-carstens/golastic"
)

type indexSettings struct {
	filterAction
	options *options.IndexSettingsOptions
}

func (is *indexSettings) ApplyOptions() Actionable {
	is.options = is.context.Options().(*options.IndexSettingsOptions)

	is.builder.SetOptions(&golastic.IndexOptions{Timeout: is.options.TimeoutInSeconds()})

	return is
}

func (is *indexSettings) Perform() Actionable {
	is.exec(func(index string) error {
		settings, err := mapToString(is.options.IndexSettings)

		if err != nil {
			return err
		}

		response, err := is.builder.PutSettings(settings, index)

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
