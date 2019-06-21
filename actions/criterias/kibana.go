package criterias

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
)

type Kibana struct {
	baseCriteria
}

func (k *Kibana) Validate() error {
	return nil
}

func (k *Kibana) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	err := json.Unmarshal(container.Bytes(), k)

	return k, err
}

func (k *Kibana) Name() string {
	return "kibana"
}
