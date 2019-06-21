package criterias

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
)

type Count struct {
	baseUseAge
	Count   int    `json:"count"`
	Pattern string `json:"pattern"`
}

func (c *Count) Validate() error {
	if c.Count < 1 {
		return errors.New("Count needs to be greater than 0.")
	}

	return c.validateUseAge()
}

func (c *Count) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	err := json.Unmarshal(container.Bytes(), c)

	return c, err
}

func (c *Count) Name() string {
	return "count"
}
