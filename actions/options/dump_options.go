package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

const DEFAULT_DUMP_MAX_EXECUTION_TIME int = 1
const DEFAULT_DUMP_CONCURRENCY int = 5
const FS_REPOSITORY string = "fs"
const GCS_REPOSITORY string = "gcs"

type DumpOptions struct {
	defaultOptions
	Criteria            []*QueryCriteria `json:"criteria"`
	MaxExecutionTime    int              `json:"max_execution_time"`
	Repository          string           `json:"repository"`
	Name                string           `json:"name"`
	Path                string           `json:"path"`
	Bucket              string           `json:"bucket"`
	CredentialsFilePath string           `json:"credentials_file_path"`
	Concurrency         int              `json:"concurrency"`
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

	return nil
}

func (do *DumpOptions) Validate() error {
	if do.MaxExecutionTime < 0 {
		return errors.New("max_execution_time must be greater than or equal to 0")
	}

	if len(do.Name) == 0 {
		return errors.New("name cannot be empty")
	}

	if len(do.Path) == 0 {
		return errors.New("path cannot be empty")
	}

	if !inStringSlice(do.Repository, []string{FS_REPOSITORY, GCS_REPOSITORY}) {
		return errors.New("invalid repsitory, please select either fs or gcs")
	}

	if do.Concurrency > 10 {
		return errors.New("concurrency cannot be greater than 10")
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
