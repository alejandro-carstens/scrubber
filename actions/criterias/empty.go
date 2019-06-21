package criterias

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
)

type Empty struct {
	baseCriteria
}

func (e *Empty) Validate() error {
	return nil
}

func (e *Empty) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	err := json.Unmarshal(container.Bytes(), e)

	return e, err
}

func (e *Empty) Name() string {
	return "empty"
}
