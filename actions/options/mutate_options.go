package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type MutateOptions struct {
	defaultOptions
	Criteria           []*QueryCriteria       `json:"criteria"`
	BatchSize          int                    `json:"batch_size"`
	MaxExecutionTime   int                    `json:"max_execution_time"`
	RetryCountPerQuery int                    `json:"retry_count_per_query"`
	Action             string                 `json:"action"`
	Mutation           map[string]interface{} `json:"mutation"`
}

func (mo *MutateOptions) FillFromContainer(container *gabs.Container) error {
	mo.container = container

	return json.Unmarshal(container.Bytes(), mo)
}

func (mo *MutateOptions) Validate() error {
	for _, criteria := range mo.Criteria {
		if err := criteria.validate(); err != nil {
			return err
		}
	}

	if !inStringSlice(mo.Action, []string{"update", "delete"}) {
		return errors.New("invalid action, please select either insert, update or delete]")
	}

	if mo.Mutation == nil && mo.Action == "update" {
		return errors.New("you must specify a mutation for the update action")
	}

	if mo.BatchSize < 0 {
		return errors.New("batch_size can't be negative")
	}

	if mo.MaxExecutionTime < 0 {
		return errors.New("max_execution_time can't be negative")
	}

	if mo.RetryCountPerQuery < 0 {
		return errors.New("retry_count_per_query can't be negative")
	}

	return nil
}

func (so *MutateOptions) BindFlags(flags *pflag.FlagSet) error {
	return nil
}
