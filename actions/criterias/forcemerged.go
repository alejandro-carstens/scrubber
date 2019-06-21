package criterias

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
)

type Forcemerged struct {
	baseCriteria
	MaxNumSegments int `json:"max_segments"`
}

func (f *Forcemerged) Validate() error {
	if f.MaxNumSegments < 0 {
		return errors.New("max_num_segments should be greater than 0")
	}

	return nil
}

func (f *Forcemerged) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	err := json.Unmarshal(container.Bytes(), f)

	return f, err
}

func (f *Forcemerged) Name() string {
	return "forcemerged"
}
