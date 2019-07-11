package actions

import (
	"errors"
	"scrubber/actions/options"

	"github.com/alejandro-carstens/golastic"
)

const DEFAULT_MAX_SNAPSHOT_BYTES_PER_SECOND string = "20mb"
const DEFAULT_MAX_RESTORE_BYTES_PER_SECOND string = "20mb"
const DEFAULT_COMPRESS bool = true
const DEFAULT_VERIFY bool = true

type createRepository struct {
	action
	options *options.CreateRepositoryOptions
}

func (cr *createRepository) ApplyOptions() Actionable {
	cr.options = cr.context.Options().(*options.CreateRepositoryOptions)

	if !cr.options.Exists("compress") {
		cr.options.Compress = DEFAULT_COMPRESS
	}

	if !cr.options.Exists("verify") {
		cr.options.Verify = DEFAULT_VERIFY
	}

	if !cr.options.Exists("max_restore_bytes_per_second") {
		cr.options.MaxRestoreBytesPerSecond = DEFAULT_MAX_RESTORE_BYTES_PER_SECOND
	}

	if !cr.options.Exists("max_snapshot_bytes_per_second") {
		cr.options.MaxSnapshotBytesPerSecond = DEFAULT_MAX_SNAPSHOT_BYTES_PER_SECOND
	}

	cr.builder.SetOptions(&golastic.IndexOptions{Timeout: cr.options.TimeoutInSeconds()})

	return cr
}

func (cr *createRepository) Perform() Actionable {
	settings := map[string]interface{}{}

	settings["compress"] = cr.options.Compress
	settings["max_restore_bytes_per_sec"] = cr.options.MaxRestoreBytesPerSecond
	settings["max_snapshot_bytes_per_sec"] = cr.options.MaxSnapshotBytesPerSecond

	if cr.options.Exists("chunk_size") {
		settings["chunk_size"] = cr.options.ChunkSize
	}

	if cr.options.RepoType == "fs" {
		settings["location"] = cr.options.Location
	} else if cr.options.RepoType == "gcs" {
		settings["bucket"] = cr.options.Bucket
	}

	response, err := cr.builder.CreateRepository(cr.options.Repository, cr.options.RepoType, cr.options.Verify, settings)

	if err != nil {
		cr.errorReportMap.push(cr.name, "_all", err)

		return cr
	}

	if acknowledged, _ := response.S("acknowledged").Data().(bool); !acknowledged {
		cr.errorReportMap.push(cr.name, "_all", errors.New("repository was not acknowledge"))
	}

	return cr
}

func (cr *createRepository) ApplyFilters() error {
	return nil
}
