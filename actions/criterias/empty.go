package criterias

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
)

// Empty representation of the empty filter criteria
type Empty struct {
	baseCriteria
}

// Validate implementation of the Criteriable interface
func (e *Empty) Validate() error {
	return nil
}

// FillFromContainer implementation of the Criteriable interface
func (e *Empty) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	err := json.Unmarshal(container.Bytes(), e)

	return e, err
}

// Name implementation of the Criteriable interface
func (e *Empty) Name() string {
	return "empty"
}
