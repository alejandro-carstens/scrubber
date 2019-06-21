package criterias

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
)

type Closed struct {
	baseCriteria
}

func (c *Closed) Validate() error {
	return nil
}

func (c *Closed) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	err := json.Unmarshal(container.Bytes(), c)

	return c, err
}

func (c *Closed) Name() string {
	return "closed"
}
