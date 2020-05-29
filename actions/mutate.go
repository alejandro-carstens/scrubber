package actions

import (
	"errors"
	"time"

	"scrubber/actions/options"

	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/golastic"
)

type mutate struct {
	filterAction
	options *options.MutateOptions
}

// ApplyOptions implementation of the Actionable interface
func (m *mutate) ApplyOptions() Actionable {
	m.options = m.context.Options().(*options.MutateOptions)

	return m
}

// Perform implementation of the Actionable interface
func (m *mutate) Perform() Actionable {
	m.exec(func(index string) error {
		builder := m.buildQuery(index)

		if m.options.Action == "update" {
			return m.update(builder)
		}

		return m.delete(builder)
	})

	return m
}

func (m *mutate) buildQuery(index string) *golastic.Builder {
	builder := m.connection.Builder(index)
	buildQuery(builder, m.options.Criteria)

	return builder
}

func (m *mutate) update(builder *golastic.Builder) error {
	taskResponse, err := builder.ExecuteAsync(m.options.BatchSize, m.options.Mutation)

	if err != nil {
		return err
	}

	m.reporter.logger.Noticef("Response: %v", taskResponse.String())

	return m.monitorTask(taskResponse, builder)
}

func (m *mutate) delete(builder *golastic.Builder) error {
	taskResponse, err := builder.DestroyAsync()

	if err != nil {
		return err
	}

	m.reporter.logger.Noticef("Response: %v", taskResponse.String())

	return m.monitorTask(taskResponse, builder)
}

func (m *mutate) monitorTask(taskResponse *gabs.Container, builder *golastic.Builder) error {
	taskId, valid := taskResponse.S("task").Data().(string)

	if !valid {
		return errors.New("could not parse the task_id")
	}

	timer := new(timer).start(int64(m.options.MaxExecutionTime))

	for {
		task, err := builder.GetTask(taskId, false)

		if err != nil {
			return err
		}

		m.reporter.logger.Debugf(task.String())

		completed, valid := task.S("completed").Data().(bool)

		if !valid {
			return errors.New("could not parse task completed field")
		}

		if completed {
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
