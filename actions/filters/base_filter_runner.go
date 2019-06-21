package filters

import (
	"errors"
	"scrubber/actions/criterias"
	"scrubber/actions/filters/runners/reports"

	"github.com/alejandro-carstens/golastic"
)

type baseFilterRunner struct {
	builder              *criterias.Builder
	elasticsearchBuilder *golastic.ElasticsearchBuilder
	reports              []reports.Reportable
}

func (bfr *baseFilterRunner) Init(builder *criterias.Builder, elasticsearchBuilder *golastic.ElasticsearchBuilder) error {
	if builder == nil {
		return errors.New("Builder can't be nil")
	}

	bfr.builder = builder

	if elasticsearchBuilder == nil {
		elasticsearchBuilder, err := golastic.NewBuilder(nil, nil)

		if err != nil {
			return err
		}

		bfr.elasticsearchBuilder = elasticsearchBuilder

		return nil
	}

	bfr.elasticsearchBuilder = elasticsearchBuilder
	bfr.reports = []reports.Reportable{}

	return nil
}

func (bfr *baseFilterRunner) AddReport(report reports.Reportable) {
	bfr.reports = append(bfr.reports, report)
}

func (bfr *baseFilterRunner) GetReports() []reports.Reportable {
	return bfr.reports
}
