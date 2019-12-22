package criterias

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
)

// Count represents the count filter criteria
type Count struct {
	baseUseAge
	Count   int    `json:"count"`
	Pattern string `json:"pattern"`
}

// Validate implementation of the Criteriable interface
func (c *Count) Validate() error {
	if c.Count < 1 {
		return errors.New("Count needs to be greater than 0.")
	}

	return c.validateUseAge()
}

// FillFromContainer implementation of the Criteriable interface
func (c *Count) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	err := json.Unmarshal(container.Bytes(), c)

	return c, err
}

// Name implementation of the Criteriable interface
func (c *Count) Name() string {
	return "count"
}
