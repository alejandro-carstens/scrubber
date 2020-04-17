package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type ImportOptions struct {
	defaultOptions
	MaxExecutionTime    int    `json:"max_execution_time"`
	Repository          string `json:"repository"`
	Name                string `json:"name"`
	Path                string `json:"path"`
	Bucket              string `json:"bucket"`
	CredentialsFilePath string `json:"credentials_file_path"`
	Concurrency         int    `json:"concurrency"`
	Size                int    `json:"size"`
}

func (io *ImportOptions) FillFromContainer(container *gabs.Container) error {
	io.container = container

	if err := json.Unmarshal(container.Bytes(), io); err != nil {
		return err
	}

	if io.MaxExecutionTime == 0 {
		io.MaxExecutionTime = DEFAULT_DUMP_MAX_EXECUTION_TIME
	}

	if io.Concurrency <= 0 {
		io.Concurrency = DEFAULT_DUMP_CONCURRENCY
	}

	if len(io.Repository) == 0 {
		io.Repository = FS_REPOSITORY
	}

	if io.Size == 0 {
		io.Size = 2500
	}

	return nil
}

func (io *ImportOptions) Validate() error {
	if io.MaxExecutionTime < 0 {
		return errors.New("max_execution_time must be greater than or equal to 0")
	}

	if len(io.Name) == 0 {
		return errors.New("name cannot be empty")
	}

	if !inStringSlice(io.Repository, []string{FS_REPOSITORY, GCS_REPOSITORY}) {
		return errors.New("invalid repsitory, please select either fs or gcs")
	}

	if io.Repository == FS_REPOSITORY && len(io.Path) == 0 {
		return errors.New("path cannot be empty")
	}

	if io.Concurrency > 10 {
		return errors.New("concurrency cannot be greater than 10")
	}

	if io.Size <= 0 || io.Size > 10000 {
		return errors.New("size must be greater than 0 and lesser than or equal to 10000")
	}

	if io.Repository == "gcs" {
		if len(io.Bucket) == 0 {
			return errors.New("a bucket needs to be specified when using the gcs repository")
		}

		if len(io.CredentialsFilePath) == 0 {
			return errors.New("a credentials_file_path needs to be specified when using the gcs repository")
		}
	}

	return nil
}

func (io *ImportOptions) BindFlags(flags *pflag.FlagSet) error {
	return nil
}
