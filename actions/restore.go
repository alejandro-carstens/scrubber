package actions

import (
	"errors"
	"regexp"
	"scrubber/actions/options"
	"scrubber/actions/infos"
	"strings"
	"time"

	"github.com/alejandro-carstens/golastic"
)

type restore struct {
	snapshotAction
	options *options.RestoreOptions
}

func (r *restore) ApplyOptions() Actionable {
	r.options = r.context.Options().(*options.RestoreOptions)

	if !r.options.Exists("wait_for_completion") {
		r.options.WaitForCompletion = DEFAULT_WAIT_FOR_COMPLETION
	}

	if !r.options.WaitForCompletion {
		if !r.options.Exists("max_wait") {
			r.options.MaxWait = DEFAULT_MAX_WAIT
		}

		if !r.options.Exists("wait_interval") {
			r.options.WaitInterval = DEFAULT_WAIT_INTERVAL
		}
	}

	r.builder.SetOptions(&golastic.IndexOptions{
		Timeout:            r.options.TimeoutInSeconds(),
		WaitForCompletion:  r.options.WaitForCompletion,
		Partial:            r.options.Partial,
		IncludeGlobalState: r.options.IncludeGlobalState,
		IgnoreUnavailable:  r.options.IgnoreUnavailable,
		IncludeAliases:     r.options.IncludeAliases,
		Indices:            strings.Split(r.options.Indices, ","),
		RenamePattern:      r.options.RenamePattern,
		RenameReplacement:  r.options.RenameReplacement,
		IndexSettings:      r.options.ExtraSettings,
	})

	return r
}

func (r *restore) Perform() Actionable {
	if len(r.list) == 0 {
		return r
	}

	snapshot := r.list[0]

	snapshotsInProgress, err := r.checkForSnapshotsInProgress(r.options.Repository)

	if err != nil {
		r.errorReportMap.push(r.name, snapshot, err)

		return r
	}

	if snapshotsInProgress {
		r.errorReportMap.push(r.name, snapshot, errors.New("Snapshot currently in progress"))

		return r
	}

	if r.options.WaitForCompletion {
		if err := r.runAndWaitForCompletion(snapshot); err != nil {
			r.errorReportMap.push(r.name, snapshot, err)
		}

		return r
	}

	response, err := r.builder.SnapshotRestore(r.options.Repository, snapshot)

	if err != nil {
		r.errorReportMap.push(r.name, snapshot, err)

		return r
	}

	if accepted, _ := response.S("accepted").Data().(bool); !accepted {
		r.errorReportMap.push(r.name, snapshot, errors.New("Restore was not accepted"))

		return r
	}

	if err := r.checkRestoreStatus(snapshot); err != nil {
		r.errorReportMap.push(r.name, snapshot, err)
	}

	return r
}

func (r *restore) runAndWaitForCompletion(snapshot string) error {
	response, err := r.builder.SnapshotRestore(r.options.Repository, snapshot)

	if err != nil {
		return err
	}

	indices, err := response.S("snapshot", "indices").Children()

	if err != nil {
		return err
	}

	if len(indices) != len(r.info[snapshot].(*infos.SnapshotInfo).Indices) {
		return errors.New("One or more indices failed to restore")
	}

	return nil
}

func (r *restore) checkRestoreStatus(snapshot string) error {
	indexList := strings.Split(r.options.Indices, ",")

	if len(indexList) == 0 {
		for _, index := range r.info[snapshot].(*infos.SnapshotInfo).Indices {
			indexList = append(indexList, index)
		}
	}

	expectedIndexList := []string{}

	if len(r.options.RenamePattern) == 0 || len(r.options.RenameReplacement) == 0 {
		expectedIndexList = indexList
	} else {
		for _, index := range indexList {
			expectedIndexList = append(expectedIndexList, r.formatExpectedIndex(index))
		}
	}

	counter := 0

	for {
		if counter >= r.options.MaxWait {
			return errors.New("Max wait time reached")
		}

		recovered, err := r.checkRecoveryStatus(expectedIndexList...)

		if err != nil {
			return err
		}

		if recovered {
			return nil
		}

		counter = counter + r.options.WaitInterval

		time.Sleep(time.Duration(int64(r.options.WaitInterval)) * time.Second)
	}
}

func (r *restore) checkRecoveryStatus(indices ...string) (bool, error) {
	recovery, err := r.builder.Recovery(indices...)

	if err != nil {
		return false, err
	}

	if len(recovery) == 0 {
		return false, nil
	}

	for _, index := range recovery {
		shards, err := index.S("shards").Children()

		if err != nil {
			return false, err
		}

		for _, shard := range shards {
			stage, valid := shard.S("stage").Data().(string)

			if !valid {
				return false, errors.New("Could not retrieve stage from shard response")
			}

			if stage != "DONE" {
				return false, nil
			}
		}
	}

	return true, nil
}

func (r *restore) formatExpectedIndex(index string) string {
	return regexp.MustCompile(r.options.RenamePattern).ReplaceAllString(index, r.options.RenameReplacement)
}
