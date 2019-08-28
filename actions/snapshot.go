package actions

import (
	"errors"
	"fmt"
	"scrubber/actions/options"
	"time"

	"github.com/alejandro-carstens/golastic"
)

const DEFAULT_WAIT_FOR_COMPLETION bool = true
const DEFAULT_MAX_WAIT int = 3600
const DEFAULT_WAIT_INTERVAL int = 9
const SUCCESS_STATUS string = "SUCCESS"
const FAILED_STATUS string = "FAILED"
const PARTIAL_STATUS string = "PARTIAL"
const IN_PROGRESS_STATUS string = "IN_PROGRESS"

type snapshot struct {
	snapshotAction
	options *options.SnapshotOptions
}

func (s *snapshot) ApplyOptions() Actionable {
	s.options = s.context.Options().(*options.SnapshotOptions)

	if len(s.options.Name) == 0 {
		s.options.Name = "scrubber-" + time.Now().Format("1992-06-02")
	}

	if !s.options.Exists("wait_for_completion") {
		s.options.WaitForCompletion = DEFAULT_WAIT_FOR_COMPLETION
	}

	if !s.options.WaitForCompletion {
		if !s.options.Exists("max_wait") {
			s.options.MaxWait = DEFAULT_MAX_WAIT
		}

		if !s.options.Exists("wait_interval") {
			s.options.WaitInterval = DEFAULT_WAIT_INTERVAL
		}
	}

	s.indexer.SetOptions(&golastic.IndexOptions{
		Timeout:            s.options.TimeoutInSeconds(),
		WaitForCompletion:  s.options.WaitForCompletion,
		Partial:            s.options.Partial,
		IncludeGlobalState: s.options.IncludeGlobalState,
		IgnoreUnavailable:  s.options.IgnoreUnavailable,
	})

	return s
}

func (s *snapshot) Perform() Actionable {
	snapshotsInProgress, err := s.checkForSnapshotsInProgress(s.options.Repository)

	if err != nil {
		return s.logError(err)
	}

	if snapshotsInProgress {
		return s.logError(errors.New("There is a snapshot currently in progress"))
	}

	if s.options.WaitForCompletion {
		if err := s.runAndWaitForCompletion(); err != nil {
			return s.logError(err)
		}

		return s
	}

	response, err := s.indexer.Snapshot(s.options.Repository, s.options.Name, s.list...)

	if err != nil {
		return s.logError(err)
	}

	if accepted, _ := response.S("accepted").Data().(bool); !accepted {
		return s.logError(errors.New("Snapshot was not accepted"))
	}

	if err := s.checkActionStatus(); err != nil {
		return s.logError(err)
	}

	return s
}

func (s *snapshot) runAndWaitForCompletion() error {
	response, err := s.indexer.Snapshot(s.options.Repository, s.options.Name, s.list...)

	if err != nil {
		return err
	}

	if state, valid := response.S("snapshot", "state").Data().(string); !valid || state == FAILED_STATUS {
		return errors.New("Snapshot was not successful")
	}

	return nil
}

func (s *snapshot) checkActionStatus() error {
	counter := 0

	for {
		if counter >= s.options.MaxWait {
			return errors.New("max_wait reached")
		}

		response, err := s.indexer.GetSnapshots(s.options.Repository, s.options.Name)

		if err != nil {
			return err
		}

		snapshots, err := response.S("snapshots").Children()

		if err != nil {
			return err
		}

		if len(snapshots) > 1 {
			return errors.New("There should only be one snapshot in progress")
		}

		state, valid := snapshots[0].S("state").Data().(string)

		if !valid {
			return errors.New("Could not retrieve state from snapshot response")
		}

		switch state {
		case SUCCESS_STATUS:
			return nil
		case PARTIAL_STATUS:
			return nil
		case FAILED_STATUS:
			return errors.New("Snapshot was not successful")
		}

		counter = counter + s.options.WaitInterval

		time.Sleep(time.Duration(int64(s.options.WaitInterval)) * time.Second)
	}
}

func (s *snapshot) logError(err error) *snapshot {
	s.errorReportMap.push(s.name, fmt.Sprint(s.list), err)

	return s
}
