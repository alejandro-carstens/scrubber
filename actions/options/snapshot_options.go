package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type SnapshotOptions struct {
	baseSnapshotOptions
	Name               string `json:"name"`
	IgnoreUnavailable  bool   `json:"ignore_unavailable"`
	IncludeGlobalState bool   `json:"include_global_state"`
	Partial            bool   `json:"partial"`
	WaitForCompletion  bool   `json:"wait_for_completion"`
	MaxWait            int    `json:"max_wait"`
	WaitInterval       int    `json:"wait_interval"`
}

func (so *SnapshotOptions) FillFromContainer(container *gabs.Container) error {
	so.container = container

	return json.Unmarshal(container.Bytes(), so)
}

func (so *SnapshotOptions) Validate() error {
	if len(so.Repository) == 0 {
		return errors.New("repository value is required")
	}

	return nil
}

func (so *SnapshotOptions) BindFlags(flags *pflag.FlagSet) error {
	so.defaultBindFlags(flags)

	name, _ := flags.GetString("name")
	repository, _ := flags.GetString("repository")
	ignoreUnavailable, _ := flags.GetBool("ignore_unavailable")
	includeGlobalState, _ := flags.GetBool("include_global_state")
	partial, _ := flags.GetBool("partial")
	waitForCompletion, _ := flags.GetBool("wait_for_completion")
	maxWait, _ := flags.GetInt("max_wait")
	waitInterval, _ := flags.GetInt("wait_interval")

	so.Repository = repository
	so.Name = name
	so.IgnoreUnavailable = ignoreUnavailable
	so.IncludeGlobalState = includeGlobalState
	so.Partial = partial
	so.WaitForCompletion = waitForCompletion
	so.MaxWait = maxWait
	so.WaitInterval = waitInterval

	return nil
}
