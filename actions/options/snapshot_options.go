package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type SnapshotOptions struct {
	defaultOptions
	Repository         string `json:"repository"`
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

	so.Repository = stringFromFlags(flags, "repository")
	so.Name = stringFromFlags(flags, "name")
	so.IgnoreUnavailable = boolFromFlags(flags, "ignore_unavailable")
	so.IncludeGlobalState = boolFromFlags(flags, "include_global_state")
	so.Partial = boolFromFlags(flags, "partial")
	so.WaitForCompletion = boolFromFlags(flags, "wait_for_completion")
	so.MaxWait = intFromFlags(flags, "max_wait")
	so.WaitInterval = intFromFlags(flags, "wait_interval")

	return nil
}
