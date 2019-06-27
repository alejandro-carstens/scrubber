package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type CreateRepositoryOptions struct {
	defaultOptions
	Compress                  bool   `json:"compress"`
	ChunkSize                 string `json:"chunk_size"`
	MaxRestoreBytesPerSecond  string `json:"max_restore_bytes_per_second"`
	MaxSnapshotBytesPerSecond string `json:"max_snapshot_bytes_per_second"`
	Location                  string `json:"location"`
	Repository                string `json:"repository"`
	RepoType                  string `json:"repo_type"`
	Verify                    bool   `json:"verify"`
}

func (cro *CreateRepositoryOptions) FillFromContainer(container *gabs.Container) error {
	cro.container = container

	return json.Unmarshal(container.Bytes(), cro)
}

func (cro *CreateRepositoryOptions) Validate() error {
	if len(cro.RepoType) == 0 {
		return errors.New("repo_type is a required option")
	}

	if cro.RepoType != "fs" {
		return errors.New("repo_type must be of type 'fs'")
	}

	if len(cro.Location) == 0 {
		return errors.New("location option is required")
	}

	if len(cro.Repository) == 0 {
		return errors.New("repository option is required")
	}

	return nil
}

func (cro *CreateRepositoryOptions) BindFlags(flags *pflag.FlagSet) error {
	cro.defaultBindFlags(flags)

	compress, _ := flags.GetBool("compress")
	chunkSize, _ := flags.GetString("chunk_size")
	maxRestoreBytesPerSecond, _ := flags.GetString("max_restore_bytes_per_second")
	maxSnapshotBytesPerSecond, _ := flags.GetString("max_snapshot_bytes_per_second")
	location, _ := flags.GetString("location")
	repository, _ := flags.GetString("repository")
	repoType, _ := flags.GetString("repo_type")
	verify, _ := flags.GetBool("verify")

	cro.Compress = compress
	cro.ChunkSize = chunkSize
	cro.MaxRestoreBytesPerSecond = maxRestoreBytesPerSecond
	cro.MaxSnapshotBytesPerSecond = maxSnapshotBytesPerSecond
	cro.Location = location
	cro.Repository = repository
	cro.RepoType = repoType
	cro.Verify = verify

	return nil
}
