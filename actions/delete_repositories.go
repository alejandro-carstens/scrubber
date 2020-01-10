package actions

import (
	"errors"
	"strings"

	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

type deleteRepositories struct {
	action
	options *options.DeleteRepositoriesOptions
}

// ApplyOptions implementation of the actionable interface
func (dr *deleteRepositories) ApplyOptions() Actionable {
	dr.options = dr.context.Options().(*options.DeleteRepositoriesOptions)

	dr.indexer.SetOptions(&golastic.IndexOptions{Timeout: dr.options.TimeoutInSeconds()})

	return dr
}

// Perform implementation of the actionable interface
func (dr *deleteRepositories) Perform() Actionable {
	response, err := dr.indexer.DeleteRepositories(strings.Split(dr.options.Repositories, ",")...)

	if err != nil {
		dr.errorContainer.push(dr.name, dr.options.Repositories, err)

		return dr
	}

	if acknowledged, _ := response.S("acknowledged").Data().(bool); !acknowledged {
		dr.errorContainer.push(
			dr.name,
			dr.options.Repositories,
			errors.New("delete repositopries action not acknoledged"),
		)
	}

	return dr
}

// ApplyFilters implementation of the Actionable interface
func (dr *deleteRepositories) ApplyFilters() error {
	return nil
}
