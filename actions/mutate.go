package actions

import (
	"errors"
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
	response, err := builder.ExecuteAsync(m.options.BatchSize, m.options.Mutation)

	if err != nil {
		return err
	}

	m.reporter.logger.Noticef("Response: %v", response.String())

	taskId, valid := response.S("task").Data().(string)

	if !valid {
		return errors.New("could not parse the task_id")
	}

	timer := new(timer).start(int64(m.options.MaxExecutionTime))

	for {
		task, err := builder.GetTask(taskId)

		if err != nil {
			return err
		}

		m.reporter.logger.Debugf(task.String())

		complete, valid := task.S("completed").Data().(bool)

		if !valid {
			return errors.New("could not parse task completed field")
		}

		if complete {
			return nil
		}

		time.Sleep(time.Duration(int64(m.options.WaitInterval)) * time.Second)

		if timer.expired() {
			m.reporter.logger.Noticef("max_execution_time %v exceeded", m.options.MaxExecutionTime)

			break
		}
	}

	cancelledTask, err := builder.CancelTask(taskId)

	if err != nil {
		return err
	}

	m.reporter.logger.Debugf(cancelledTask.String())

	return nil
}
