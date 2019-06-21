package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
)

type RestoreOptions struct {
	baseSnapshotOptions
	IncludeAliases     bool                   `json:"include_aliases"`
	Partial            bool                   `json:"partial"`
	WaitForCompletion  bool                   `json:"wait_for_completion"`
	MaxWait            int                    `json:"max_wait"`
	WaitInterval       int                    `json:"wait_interval"`
	RenamePattern      string                 `json:"rename_pattern"`
	RenameReplacement  string                 `json:"rename_replacement"`
	Indices            string                 `json:"indices"`
	ExtraSettings      map[string]interface{} `json:"extra_settings"`
	IncludeGlobalState bool                   `json:"include_global_state"`
	IgnoreUnavailable  bool                   `json:"ignore_unavailable"`
	Name               string                 `json:"name"`
}

func (ro *RestoreOptions) FillFromContainer(container *gabs.Container) error {
	ro.container = container

	return json.Unmarshal(container.Bytes(), ro)
}

func (ro *RestoreOptions) Validate() error {
	if len(ro.Repository) == 0 {
		return errors.New("repository value is required")
	}

	if len(ro.RenamePattern) > 0 && len(ro.RenameReplacement) == 0 {
		return errors.New("rename_replacement is required if the rename_pattern option is specified")
	}

	if len(ro.RenamePattern) == 0 && len(ro.RenameReplacement) > 0 {
		return errors.New("rename_pattern is required if the rename_replacement option is specified")
	}

	return nil
}
