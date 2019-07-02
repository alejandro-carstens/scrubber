package options

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
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

func (ro *RestoreOptions) BindFlags(flags *pflag.FlagSet) error {
	ro.defaultBindFlags(flags)

	repository, _ := flags.GetString("repository")
	includeAliases, _ := flags.GetBool("include_aliases")
	partial, _ := flags.GetBool("partial")
	waitForCompletion, _ := flags.GetBool("wait_for_completion")
	includeGlobalState, _ := flags.GetBool("include_global_state")
	ignoreUnavailable, _ := flags.GetBool("ignore_unavailable")
	maxWait, _ := flags.GetInt("max_wait")
	waitInterval, _ := flags.GetInt("wait_interval")
	renamePattern, _ := flags.GetString("rename_pattern")
	renameReplacement, _ := flags.GetString("rename_replacement")
	indices, _ := flags.GetString("indices")
	name, _ := flags.GetString("name")
	extraSettings, _ := flags.GetString("extra_settings")

	if len(extraSettings) > 0 {
		es := map[string]interface{}{}

		if err := json.Unmarshal([]byte(extraSettings), &es); err != nil {
			return err
		}

		ro.ExtraSettings = es
	}

	ro.Repository = repository
	ro.IncludeAliases = includeAliases
	ro.Partial = partial
	ro.WaitForCompletion = waitForCompletion
	ro.MaxWait = maxWait
	ro.WaitInterval = waitInterval
	ro.RenamePattern = renamePattern
	ro.RenameReplacement = renameReplacement
	ro.Indices = indices
	ro.IncludeGlobalState = includeGlobalState
	ro.IgnoreUnavailable = ignoreUnavailable
	ro.Name = name

	return nil
}

func (ro *RestoreOptions) Exists(value string) bool {
	if ro.container == nil {
		ro.container = toContainer(ro)
	}

	return ro.container.Exists(value)
}

func (ro *RestoreOptions) String(value string) string {
	if ro.container == nil {
		ro.container = toContainer(ro)
	}

	return fmt.Sprint(ro.container.S(value).Data())
}
