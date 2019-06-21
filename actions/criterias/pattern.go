package criterias

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
)

type Pattern struct {
	baseCriteria
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

func (p *Pattern) Validate() error {
	if len(p.Kind) == 0 {
		return errors.New("Kind cannot be empty")
	}

	if len(p.Value) == 0 {
		return errors.New("Value cannot be empty")
	}

	if p.Kind != "suffix" && p.Kind != "prefix" && p.Kind != "regex" && p.Kind != "timestring" {
		return errors.New("Invalid kind")
	}

	if p.Kind == "timestring" {
		if err := validateTimestring(p.Value); err != nil {
			return err
		}
	}

	return nil
}

func (p *Pattern) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	err := json.Unmarshal(container.Bytes(), p)

	return p, err
}

func (p *Pattern) Name() string {
	return "pattern"
}
