package criterias

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
)

type State struct {
	baseCriteria
	State string `json:"state"`
}

func (s *State) Validate() error {
	if s.State != "SUCCESS" && s.State != "PARTIAL" && s.State != "FAILED" && s.State != "IN_PROGRESS" {
		return errors.New("invalid state value, please use: SUCCESS, PARTIAL, FAILED & IN_PROGRESS")
	}

	return nil
}

func (s *State) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	err := json.Unmarshal(container.Bytes(), s)

	return s, err
}

func (s *State) Name() string {
	return "state"
}
