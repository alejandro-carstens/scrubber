package criterias

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
)

// Forcemerged representst the forcemerged filter criteria
type Forcemerged struct {
	baseCriteria
	MaxNumSegments int `json:"max_segments"`
}

// Validate implementation of the Criteriable interface
func (f *Forcemerged) Validate() error {
	if f.MaxNumSegments < 0 {
		return errors.New("max_num_segments should be greater than 0")
	}

	return nil
}

// FillFromContainer implementation of the Criteriable interface
func (f *Forcemerged) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	err := json.Unmarshal(container.Bytes(), f)

	return f, err
}

// Name implementation of the Criteriable interface
func (f *Forcemerged) Name() string {
	return "forcemerged"
}
