package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type DumpOptions struct {
	defaultOptions
	Criteria            []*QueryCriteria `json:"criteria"`
	MaxExecutionTime    int              `json:"max_execution_time"`
	Repository          string           `json:"repository"`
	Name                string           `json:"name"`
	Path                string           `json:"path"`
	Bucket              string           `json:"bucket"`
	CredentialsFilePath string           `json:"credentials_file_path"`
	KeepAlive           int              `json:"keep_alive"`
	Concurrency         int              `json:"concurrency"`
	Size                int              `json:"size"`
}

func (do *DumpOptions) FillFromContainer(container *gabs.Container) error {
	do.container = container

	if err := json.Unmarshal(container.Bytes(), do); err != nil {
		return err
	}

	if do.MaxExecutionTime == 0 {
		do.MaxExecutionTime = DEFAULT_DUMP_MAX_EXECUTION_TIME
	}

	if do.Concurrency <= 0 {
		do.Concurrency = DEFAULT_DUMP_CONCURRENCY
	}

	if len(do.Repository) == 0 {
		do.Repository = FS_REPOSITORY
	}

	if do.KeepAlive <= 0 {
		do.KeepAlive = DEFAULT_DUMP_KEEP_ALIVE
	}

	if do.Size == 0 {
		do.Size = 2500
	}

	return nil
}

func (do *DumpOptions) Validate() error {
	if do.MaxExecutionTime < 0 {
		return errors.New("max_execution_time must be greater than or equal to 0")
	}

	if len(do.Name) == 0 {
		return errors.New("name cannot be empty")
	}

	if !inStringSlice(do.Repository, []string{FS_REPOSITORY, GCS_REPOSITORY}) {
		return errors.New("invalid repsitory, please select either fs or gcs")
	}

	if do.Repository == FS_REPOSITORY && len(do.Path) == 0 {
		return errors.New("path cannot be empty")
	}

	if do.Concurrency > 10 {
		return errors.New("concurrency cannot be greater than 10")
	}

	if do.Size <= 0 || do.Size > 10000 {
		return errors.New("size must be greater than 0 and lesser than or equal to 10000")
	}

	if do.KeepAlive <= 0 {
		return errors.New("keep_alive must be greater than 0")
	}

	if do.Repository == "gcs" {
		if len(do.Bucket) == 0 {
			return errors.New("a bucket needs to be specified when using the gcs repository")
		}

		if len(do.CredentialsFilePath) == 0 {
			return errors.New("a credentials_file_path needs to be specified when using the gcs repository")
		}
	}

	for _, criteria := range do.Criteria {
		if err := criteria.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (do *DumpOptions) BindFlags(flags *pflag.FlagSet) error {
	return nil
}
