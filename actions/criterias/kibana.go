package criterias

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
)

// Kibana represents the kibana filter criteria
type Kibana struct {
	baseCriteria
}

// Validate implementation of the Criteriable interface
func (k *Kibana) Validate() error {
	return nil
}

// Kibana implementation of the Criteriable interface
func (k *Kibana) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	err := json.Unmarshal(container.Bytes(), k)

	return k, err
}

// Name implementation of the Criteriable interface
func (k *Kibana) Name() string {
	return "kibana"
}
