package criterias

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
)

// Closed represents the closed filter cirteria
type Closed struct {
	baseCriteria
}

// Validate implementation of the Criteriable interface
func (c *Closed) Validate() error {
	return nil
}

// FillFromContainer implementation of the Criteriable interface
func (c *Closed) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	err := json.Unmarshal(container.Bytes(), c)

	return c, err
}

// Name implementation of the Criteriable interface
func (c *Closed) Name() string {
	return "closed"
}
