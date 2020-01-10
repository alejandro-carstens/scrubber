package actions

import (
	"errors"

	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

// DEFAULT_MAX_SNAPSHOT_BYTES_PER_SECOND is the default
// option for the max_snapshot_bytes_per_second option
const DEFAULT_MAX_SNAPSHOT_BYTES_PER_SECOND string = "20mb"

// DEFAULT_MAX_RESTORE_BYTES_PER_SECOND is the default
// option for the max_restore_bytes_per_second option
const DEFAULT_MAX_RESTORE_BYTES_PER_SECOND string = "20mb"

// DEFAULT_COMPRESS is the default
// option for the compress option
const DEFAULT_COMPRESS bool = true

// DEFAULT_VERIFY is the default
// value for the verify option
const DEFAULT_VERIFY bool = true

type createRepository struct {
	action
	options *options.CreateRepositoryOptions
}

// ApplyOptions implementation of the Actionable interface
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

	cr.indexer.SetOptions(&golastic.IndexOptions{Timeout: cr.options.TimeoutInSeconds()})

	return cr
}

// Perform implementation of the Actionable interface
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

	response, err := cr.indexer.CreateRepository(cr.options.Repository, cr.options.RepoType, cr.options.Verify, settings)

	if err != nil {
		cr.errorContainer.push(cr.name, "_all", err)

		return cr
	}

	if acknowledged, _ := response.S("acknowledged").Data().(bool); !acknowledged {
		cr.errorContainer.push(cr.name, "_all", errors.New("repository was not acknowledge"))
	}

	return cr
}

// ApplyFilters implementation of the Actionable interface
func (cr *createRepository) ApplyFilters() error {
	return nil
}
