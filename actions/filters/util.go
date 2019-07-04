package filters

import (
	"errors"
	"scrubber/actions/criterias"
	"scrubber/actions/responses"

	"github.com/alejandro-carstens/golastic"
)

func NewFilterRunner(info responses.Informable, builder *criterias.Builder, elasticsearchBuilder *golastic.ElasticsearchBuilder, isSnapshotAction bool) (*FilterRunner, error) {
	filterRunner := new(FilterRunner)

	filterRunner.info = info

	err := filterRunner.Init(builder, elasticsearchBuilder)

	return filterRunner, err
}

func NewAggregateFilterRunner(info []responses.Informable, builder *criterias.Builder, elasticsearchBuilder *golastic.ElasticsearchBuilder) (*AggregateFilterRunner, error) {
	aggregateFilterRunner := new(AggregateFilterRunner)

	if len(info) == 0 {
		return nil, errors.New("info parameter cannot be empty")
	}

	aggregateFilterRunner.info = info
	aggregateFilterRunner.countMap = map[string]int{}

	err := aggregateFilterRunner.Init(builder, elasticsearchBuilder)

	return aggregateFilterRunner, err
}
