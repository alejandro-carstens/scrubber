package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
)

type CreateRepositoryOptions struct {
	defaultOptions
	Compress                  bool   `json:"compress"`
	ChunkSize                 string `json:"chunkSize"`
	MaxRestoreBytesPerSecond  string `json:"maxRestoreBytesPerSecond"`
	MaxSnapshotBytesPerSecond string `json:"maxSnapshotBytesPerSecond"`
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
