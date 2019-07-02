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

	ro.Repository = stringFromFlags(flags, "repository")
	ro.IncludeAliases = boolFromFlags(flags, "include_aliases")
	ro.Partial = boolFromFlags(flags, "partial")
	ro.WaitForCompletion = boolFromFlags(flags, "wait_for_completion")
	ro.MaxWait = intFromFlags(flags, "max_wait")
	ro.WaitInterval = intFromFlags(flags, "wait_interval")
	ro.RenamePattern = stringFromFlags(flags, "rename_pattern")
	ro.RenameReplacement = stringFromFlags(flags, "rename_replacement")
	ro.Indices = stringFromFlags(flags, "indices")
	ro.IncludeGlobalState = boolFromFlags(flags, "include_global_state")
	ro.IgnoreUnavailable = boolFromFlags(flags, "ignore_unavailable")
	ro.Name = stringFromFlags(flags, "name")

	if len(stringFromFlags(flags, "extra_settings")) > 0 {
		ro.ExtraSettings = map[string]interface{}{}

		if err := json.Unmarshal([]byte(stringFromFlags(flags, "extra_settings")), &ro.ExtraSettings); err != nil {
			return err
		}
	}

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
