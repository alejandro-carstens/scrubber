package criterias

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
)

type Alias struct {
	baseCriteria
	Aliases []string `json:"aliases"`
}

func (a *Alias) Name() string {
	return "alias"
}

func (a *Alias) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	err := json.Unmarshal(container.Bytes(), a)

	return a, err
}

func (a *Alias) Validate() error {
	if len(a.Aliases) == 0 {
		return errors.New("Aliases cannot be empty")
	}

	return nil
}
