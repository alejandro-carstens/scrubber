package criterias

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Jeffail/gabs"
)

type Space struct {
	baseUseAge
	FeThresholdBehavior string `json:"fe_threshold_behavior"`
	DiskSpace           int    `json:"disk_space"`
	Units               string `json:"units"`
}

func (s *Space) Validate() error {
	if s.DiskSpace < 0 {
		return errors.New("disk_space needs to be a positive integer.")
	}

	if s.FeThresholdBehavior != "greater_than" && s.FeThresholdBehavior != "less_than" {
		return errors.New("fe_threshold_behavior needs to be either greater_than or less_than.")
	}

	if s.Units != "GB" && s.Units != "MB" {
		return errors.New("units needs to be either GB or MB.")
	}

	return s.validateUseAge()
}

func (s *Space) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	log.Println(container.String())

	err := json.Unmarshal(container.Bytes(), s)

	if err != nil {
		return nil, err
	}

	if len(s.Units) == 0 {
		s.Units = "GB"
	}

	if len(s.FeThresholdBehavior) == 0 {
		s.FeThresholdBehavior = "greater_than"
	}

	return s, nil
}

func (s *Space) Name() string {
	return "space"
}
