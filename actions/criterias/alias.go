package criterias

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
)

// Alias represents the alias filter criteria
type Alias struct {
	baseCriteria
	Aliases []string `json:"aliases"`
}

// Name implementation of the Criteriable interface
func (a *Alias) Name() string {
	return "alias"
}

// FillFromContainer implementation of the Criteriable interface
func (a *Alias) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	err := json.Unmarshal(container.Bytes(), a)

	return a, err
}

// Validate implementation of the Criteriable interface
func (a *Alias) Validate() error {
	if len(a.Aliases) == 0 {
		return errors.New("Aliases cannot be empty")
	}

	return nil
}
