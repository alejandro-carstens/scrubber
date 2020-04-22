package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type ImportDumpOptions struct {
	defaultOptions
	MaxExecutionTime    int                    `json:"max_execution_time"`
	Repository          string                 `json:"repository"`
	Name                string                 `json:"name"`
	Path                string                 `json:"path"`
	Bucket              string                 `json:"bucket"`
	CredentialsFilePath string                 `json:"credentials_file_path"`
	Concurrency         int                    `json:"concurrency"`
	Size                int                    `json:"size"`
	CreateIndexWaitTime int64                  `json:"create_index_wait_time"`
	RemoveSettings      []string               `json:"remove_settings"`
	RemoveMappings      []string               `json:"remove_mappings"`
	RemoveAliases       []string               `json:"remove_aliases"`
	RemoveFields        []string               `json:"remove_fields"`
	ExtraSettings       map[string]interface{} `json:"extra_settings"`
	ExtraMappings       map[string]interface{} `json:"extra_mappings"`
	ExtraAliases        map[string]interface{} `json:"extra_aliases"`
	ExtraFields         map[string]interface{} `json:"extra_fields"`
}

func (ido *ImportDumpOptions) FillFromContainer(container *gabs.Container) error {
	ido.container = container

	if err := json.Unmarshal(container.Bytes(), ido); err != nil {
		return err
	}

	if ido.MaxExecutionTime == 0 {
		ido.MaxExecutionTime = DEFAULT_DUMP_MAX_EXECUTION_TIME
	}

	if ido.Concurrency <= 0 {
		ido.Concurrency = DEFAULT_DUMP_CONCURRENCY
	}

	if len(ido.Repository) == 0 {
		ido.Repository = FS_REPOSITORY
	}

	if ido.Size == 0 {
		ido.Size = 2500
	}

	if ido.CreateIndexWaitTime == 0 {
		ido.CreateIndexWaitTime = int64(3)
	}

	if ido.ExtraAliases == nil {
		ido.ExtraAliases = map[string]interface{}{}
	}

	if ido.ExtraSettings == nil {
		ido.ExtraSettings = map[string]interface{}{}
	}

	if ido.ExtraMappings == nil {
		ido.ExtraMappings = map[string]interface{}{}
	}

	if ido.ExtraFields == nil {
		ido.ExtraFields = map[string]interface{}{}
	}

	if ido.RemoveAliases == nil {
		ido.RemoveAliases = []string{}
	}

	if ido.RemoveSettings == nil {
		ido.RemoveSettings = []string{}
	}

	if ido.RemoveMappings == nil {
		ido.RemoveMappings = []string{}
	}

	if ido.RemoveFields == nil {
		ido.RemoveFields = []string{}
	}

	return nil
}

func (ido *ImportDumpOptions) Validate() error {
	if ido.MaxExecutionTime < 0 {
		return errors.New("max_execution_time must be greater than or equal to 0")
	}

	if len(ido.Name) == 0 {
		return errors.New("name cannot be empty")
	}

	if !inStringSlice(ido.Repository, []string{FS_REPOSITORY, GCS_REPOSITORY}) {
		return errors.New("invalid repsitory, please select either fs or gcs")
	}

	if ido.Repository == FS_REPOSITORY && len(ido.Path) == 0 {
		return errors.New("path cannot be empty")
	}

	if ido.Concurrency > 10 {
		return errors.New("concurrency cannot be greater than 10")
	}

	if ido.Size <= 0 || ido.Size > 10000 {
		return errors.New("size must be greater than 0 and lesser than or equal to 10000")
	}

	if ido.Repository == "gcs" {
		if len(ido.Bucket) == 0 {
			return errors.New("a bucket needs to be specified when using the gcs repository")
		}

		if len(ido.CredentialsFilePath) == 0 {
			return errors.New("a credentials_file_path needs to be specified when using the gcs repository")
		}
	}

	return nil
}

func (ido *ImportDumpOptions) BindFlags(flags *pflag.FlagSet) error {
	return nil
}
