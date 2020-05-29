package filters

import (
	"errors"

	"scrubber/actions/criterias"
	"scrubber/actions/infos"

	"github.com/alejandro-carstens/golastic"
)

func NewFilterRunner(info infos.Informable, builder *criterias.Builder, connection *golastic.Connection) (*FilterRunner, error) {
	filterRunner := new(FilterRunner)

	filterRunner.info = info

	err := filterRunner.Init(builder, connection)

	return filterRunner, err
}

func NewAggregateFilterRunner(info []infos.Informable, builder *criterias.Builder, connection *golastic.Connection) (*AggregateFilterRunner, error) {
	aggregateFilterRunner := new(AggregateFilterRunner)

	if len(info) == 0 {
		return nil, errors.New("info parameter cannot be empty")
	}

	aggregateFilterRunner.info = info
	aggregateFilterRunner.countMap = map[string]int{}

	err := aggregateFilterRunner.Init(builder, connection)

	return aggregateFilterRunner, err
}
