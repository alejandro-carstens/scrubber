package console

import (
	"errors"
	"os"
	"path/filepath"
	"scrubber/console/tasks"
	"scrubber/ymlparser"
	"strings"

	"github.com/Jeffail/gabs"
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
}

func NewScheduler(basePath string) *Scheduler {
	scheduler := new(Scheduler)
	scheduler.basePath = basePath

	return scheduler
}

func (s *Scheduler) Run() error {
	configs, err := s.extractConfigs()

	if err != nil {
		return err
	}

	asyncActions := []*tasks.RunAction{}
	actions := []*tasks.RunAction{}
	startCron := false

	defer gocron.Clear()

	for _, config := range configs {
		task := new(tasks.RunAction).SetConfig(config)

		if config.Exists("schedule") {
			job, err := s.schedule(config)

			if err != nil {
				return err
			}

			job.Do(task.Execute)

			startCron = true

			continue
		}

		if runMode, valid := config.S("run_mode").Data().(string); valid && runMode == "async" {
			asyncActions = append(asyncActions, task)

			continue
		}

		actions = append(actions, task)
	}

	if len(actions) > 0 {
		s.runActions(actions)
	}

	if len(asyncActions) > 0 {
		s.runAsyncActions(asyncActions)
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

func (s *Scheduler) runAsyncActions(actions []*tasks.RunAction) {
	pool := grpool.NewPool(NUMBER_OF_WORKERS, len(actions))

	pool.WaitCount(len(actions))

	for _, action := range actions {
		task := action

		pool.JobQueue <- func() {
			defer pool.JobDone()

			task.Execute()
		}
	}

	pool.WaitAll()
	pool.Release()
}

func (s *Scheduler) runActions(actions []*tasks.RunAction) {
	for _, action := range actions {
		action.Execute()
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
