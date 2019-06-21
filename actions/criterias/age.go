package criterias

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
)

type Age struct {
	baseCriteria
	Source      string `json:"source"`
	Direction   string `json:"direction"`
	Units       string `json:"units"`
	UnitCount   int    `json:"unit_count"`
	Timestring  string `json:"timestring,omitempty"`
	Field       string `json:"field,omitempty"`
	StatsResult string `json:"stats_result,omitempty"`
}

func (a *Age) Name() string {
	return "age"
}

func (a *Age) FillFromContainer(container *gabs.Container) (Criteriable, error) {
	err := json.Unmarshal(container.Bytes(), a)

	return a, err
}

func (a *Age) Validate() error {
	if a.Source != "name" && a.Source != "creation_date" && a.Source != "field_stats" {
		return errors.New("Invalid Source type")
	}

	if a.Direction != "older" && a.Direction != "younger" {
		return errors.New("Invalid direction type")
	}

	if a.Units != "seconds" && a.Units != "minutes" && a.Units != "hours" &&
		a.Units != "days" && a.Units != "weeks" && a.Units != "months" && a.Units != "years" {
		return errors.New("Invalid units type")
	}

	if a.UnitCount < 1 {
		return errors.New("Unit count needs to be at least 1")
	}

	if a.Source == "name" {
		if err := validateTimestring(a.Timestring); err != nil {
			return err
		}
	}

	if a.Source == "field_stats" {
		if len(a.Field) == 0 {
			return errors.New("Since the Source is field_stats a field needs to be specified")
		}

		if a.StatsResult != "min" && a.StatsResult != "max" {
			return errors.New("Since the Source is field_stats a stats_result needs to be either min or max")
		}
	}

	return nil
}
