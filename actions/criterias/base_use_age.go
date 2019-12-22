package criterias

import (
	"errors"
)

type baseUseAge struct {
	baseCriteria
	Reverse     bool   `json:"reverse"`
	UseAge      bool   `json:"use_age"`
	Source      string `json:"source"`
	Timestring  string `json:"timestring"`
	Field       string `json:"field"`
	StatsResult string `json:"stats_result"`
	StrictMode  bool   `json:"strict_mode"`
}

// GetTimestring returns the Timestring property
func (bua *baseUseAge) GetTimestring() string {
	return bua.Timestring
}

// GetField returns the Field property
func (bua *baseUseAge) GetField() string {
	return bua.Field
}

// GetStatsResult returns the StatsResult property
func (bua *baseUseAge) GetStatsResult() string {
	return bua.StatsResult
}

// GetReverse return the Reverse property
func (bua *baseUseAge) GetReverse() bool {
	return bua.Reverse
}

// GetSource returns the Source property
func (bua *baseUseAge) GetSource() string {
	return bua.Source
}

// GetStrictMode returns the StrictMode property
func (bua *baseUseAge) GetStrictMode() bool {
	return bua.StrictMode
}

func (bua *baseUseAge) validateUseAge() error {
	if bua.UseAge {
		if len(bua.Source) == 0 {
			return errors.New("If use_age is set to true, source needs to be specified (creation_date, field_stats, or name)")
		}

		if bua.Source != "creation_date" && bua.Source != "name" && bua.Source != "field_stats" {
			return errors.New("Source needs to be either creation_date, field_stats, or name")
		}

		if bua.Source == "name" {
			return validateTimestring(bua.Timestring)
		}

		if bua.Source == "field_stats" {
			if bua.StatsResult != "min" && bua.StatsResult != "max" {
				return errors.New("Invalid stats_result field, please specify either min or max")
			}

			if len(bua.Field) == 0 {
				return errors.New("You need to specify a field for which get field_stats")
			}
		}
	}

	return nil
}
