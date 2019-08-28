package console

import (
	"errors"
	"os"
	"path/filepath"
	"scrubber/actions/contexts"
	"scrubber/logger"
	"scrubber/ymlparser"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/golastic"
	"github.com/ivpusic/grpool"
	"github.com/jasonlvhit/gocron"
)

const NUMBER_OF_WORKERS int = 50

type configMap struct {
	container *gabs.Container
	err       error
	path      string
}

type Scheduler struct {
	basePath string
	logger   *logger.Logger
	builder  *golastic.Connection
}

func (s *Scheduler) Run() error {
	configs, err := s.extractConfigs()

	if err != nil {
		return err
	}

	syncContexts := []contexts.Contextable{}
	asyncContexts := []contexts.Contextable{}
	startCron := false

	defer gocron.Clear()

	for _, config := range configs {
		context, err := contexts.New(config)

		if err != nil {
			return err
		}

		if config.Exists("schedule") {
			job, err := s.schedule(config)

			if err != nil {
				return err
			}

			job.Do(Execute, context, s.logger, s.builder)

			startCron = true

			continue
		}

		if runMode, valid := config.S("run_mode").Data().(string); valid && runMode == "async" {
			asyncContexts = append(asyncContexts, context)

			continue
		}

		syncContexts = append(syncContexts, context)
	}

	if len(syncContexts) > 0 {
		s.runActions(syncContexts)
	}

	if len(asyncContexts) > 0 {
		s.runAsyncActions(asyncContexts)
	}

	if startCron {
		<-gocron.Start()
	}

	return nil
}

func (s *Scheduler) extractConfigs() (map[string]*gabs.Container, error) {
	paths := []string{}

	if err := filepath.Walk(s.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(path, ".yml") && info.Size() > 0 {
			paths = append(paths, path)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	channel := make(chan configMap, len(paths))

	for _, path := range paths {
		filePath := path

		go func() {
			container, err := ymlparser.Parse(filePath)

			channel <- configMap{
				container: container,
				err:       err,
				path:      filePath,
			}
		}()
	}

	containers := map[string]*gabs.Container{}

	for range paths {
		confMap := <-channel

		if confMap.err != nil {
			return nil, confMap.err
		}

		containers[confMap.path] = confMap.container
	}

	return containers, nil
}

func (s *Scheduler) runAsyncActions(contexts []contexts.Contextable) {
	pool := grpool.NewPool(NUMBER_OF_WORKERS, len(contexts))

	pool.WaitCount(len(contexts))

	for _, context := range contexts {
		action := context

		pool.JobQueue <- func() {
			defer pool.JobDone()

			Execute(action, s.logger, s.builder)
		}
	}

	pool.WaitAll()
	pool.Release()
}

func (s *Scheduler) runActions(contexts []contexts.Contextable) {
	for _, context := range contexts {
		Execute(context, s.logger, s.builder)
	}
}

func (s *Scheduler) schedule(config *gabs.Container) (*gocron.Job, error) {
	value, valid := config.S("schedule", "value").Data().(float64)

	if !valid && value <= 0 {
		return nil, errors.New("Invalid or missing schedule.value on action file")
	}

	unit, valid := config.S("schedule", "unit").Data().(string)

	if !valid {
		return nil, errors.New("Invalid or missing schedule.unit on action file")
	}

	atCheck := false
	job := gocron.Every(uint64(value))

	switch strings.ToLower(unit) {
	case "seconds":
		job.Seconds()
		break
	case "minutes":
		job.Minutes()
		break
	case "hours":
		job.Hours()
		break
	case "weeks":
		job.Weeks()
		atCheck = true
		break
	case "days":
		job.Days()
		atCheck = true
		break
	case "monday":
		job.Monday()
		atCheck = true
		break
	case "tuesday":
		job.Tuesday()
		atCheck = true
		break
	case "wednesday":
		job.Wednesday()
		atCheck = true
		break
	case "thursday":
		job.Thursday()
		atCheck = true
		break
	case "friday":
		job.Friday()
		atCheck = true
		break
	case "saturday":
		job.Saturday()
		atCheck = true
		break
	case "sunday":
		job.Sunday()
		atCheck = true
		break
	default:
		return nil, errors.New("Invalid schedule.unit specified on action file")
	}

	if at, valid := config.S("schedule", "at").Data().(string); atCheck && valid {
		job.At(at)
	}

	return job, nil
}
