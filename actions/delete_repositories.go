package actions

import (
	"errors"
	"scrubber/actions/options"
	"strings"

	"github.com/alejandro-carstens/golastic"
)

type deleteRepositories struct {
	action
	options *options.DeleteRepositoriesOptions
}

func (dr *deleteRepositories) ApplyOptions() Actionable {
	dr.options = dr.context.Options().(*options.DeleteRepositoriesOptions)

	dr.indexer.SetOptions(&golastic.IndexOptions{Timeout: dr.options.TimeoutInSeconds()})

	return dr
}

func (dr *deleteRepositories) Perform() Actionable {
	response, err := dr.indexer.DeleteRepositories(strings.Split(dr.options.Repositories, ",")...)

	if err != nil {
		dr.errorReportMap.push(dr.name, dr.options.Repositories, err)

		return dr
	}

	if acknowledged, _ := response.S("acknowledged").Data().(bool); !acknowledged {
		dr.errorReportMap.push(
			dr.name,
			dr.options.Repositories,
			errors.New("delete repositopries action not acknoledged"),
		)
	}

	return dr
}

func (dr *deleteRepositories) ApplyFilters() error {
	return nil
}
