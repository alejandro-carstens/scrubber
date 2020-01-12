package actions

import (
	"time"

	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

const DEFAULT_BATCH_SIZE int = 5000

type mutate struct {
	filterAction
	options *options.MutateOptions
}

// ApplyOptions implementation of the Actionable interface
func (m *mutate) ApplyOptions() Actionable {
	m.options = m.context.Options().(*options.MutateOptions)

	m.indexer.SetOptions(&golastic.IndexOptions{Timeout: m.options.TimeoutInSeconds()})

	return m
}

// Perform implementation of the Actionable interface
func (m *mutate) Perform() Actionable {
	m.exec(func(index string) error {
		builder := m.buildQuery(index)

		if m.options.Action == "update" {
			return m.update(builder)
		}

		return nil
	})

	return m
}

func (m *mutate) buildQuery(index string) *golastic.Builder {
	builder := m.connection.Builder(index)

	buildQuery(builder, m.options.Criteria)

	if m.options.BatchSize == 0 {
		m.options.BatchSize = DEFAULT_BATCH_SIZE
	}

	builder.Limit(m.options.BatchSize)

	return builder
}

func (m *mutate) update(builder *golastic.Builder) error {
	var err error
	retryCounter := 0

	for {
		if retryCounter == m.options.RetryCountPerQuery {
			return err
		}

		count, err := builder.Count()

		if err != nil {
			retryCounter++
			m.reporter.logger.Errorf("Attempt [%v] %v", retryCounter, err.Error())

			continue
		}

		response, err := builder.Execute(m.options.Mutation)

		if response != nil {
			m.reporter.logger.Noticef("Response: %v", response.String())
		}

		if err == nil {
			if count <= int64(m.options.BatchSize) {
				break
			}

			retryCounter = 0
			time.Sleep(time.Duration(500) * time.Millisecond)

			continue
		}

		retryCounter++
		m.reporter.logger.Errorf("Attempt [%v] %v", retryCounter, err.Error())
	}

	return nil
}
