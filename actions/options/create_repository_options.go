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

	cro.Compress = boolFromFlags(flags, "compress")
	cro.ChunkSize = stringFromFlags(flags, "chunk_size")
	cro.MaxRestoreBytesPerSecond = stringFromFlags(flags, "max_restore_bytes_per_second")
	cro.MaxSnapshotBytesPerSecond = stringFromFlags(flags, "max_snapshot_bytes_per_second")
	cro.Location = stringFromFlags(flags, "location")
	cro.Repository = stringFromFlags(flags, "repository")
	cro.RepoType = stringFromFlags(flags, "repo_type")
	cro.Verify = boolFromFlags(flags, "verify")

	return nil
}
