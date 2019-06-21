package criterias

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Jeffail/gabs"
)

type Allocated struct {
	baseCriteria
	AllocationType string      `json:"allocation_type"`
	Key            string      `json:"key"`
	RawValue       interface{} `json:"value"`
	Value          string
}

func (a *Allocated) Validate() error {
	if len(a.Key) == 0 {
		return errors.New("key is a required field")
	}

	if len(a.Value) == 0 {
		return errors.New("value is a required field")
	}

	if a.AllocationType != "require" && a.AllocationType != "include" && a.AllocationType != "exclude" {
		return errors.New("allocation_type must be either: required, include, exclude or empty")
	}

	return nil
}

func (a *Allocated) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	err := json.Unmarshal(container.Bytes(), a)

	if len(a.AllocationType) == 0 {
		a.AllocationType = "require"
	}

	if a.RawValue != nil {
		a.Value = fmt.Sprint(a.RawValue)
	}

	return a, err
}

func (a *Allocated) Name() string {
	return "allocated"
}
