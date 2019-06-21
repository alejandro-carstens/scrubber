package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
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
