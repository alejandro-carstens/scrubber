package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

const DEFAULT_MAX_EXECUTION_TIME int = 1
const DEFAULT_BATCH_SIZE int = 10000

type MutateOptions struct {
	defaultOptions
	Criteria         []*QueryCriteria       `json:"criteria"`
	BatchSize        int                    `json:"batch_size"`
	MaxExecutionTime int                    `json:"max_execution_time"`
	WaitInterval     int                    `json:"wait_interval"`
	Action           string                 `json:"action"`
	Mutation         map[string]interface{} `json:"mutation"`
}

func (mo *MutateOptions) FillFromContainer(container *gabs.Container) error {
	mo.container = container

	if err := json.Unmarshal(container.Bytes(), mo); err != nil {
		return err
	}

	if mo.MaxExecutionTime == 0 {
		mo.MaxExecutionTime = DEFAULT_MAX_EXECUTION_TIME
	}

	if mo.BatchSize == 0 {
		mo.BatchSize = DEFAULT_BATCH_SIZE
	}

	return nil
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
		return errors.New("batch_size must be greater than 0 or equal to 0")
	}

	if mo.MaxExecutionTime < 0 {
		return errors.New("max_execution_time must be greater than or equal to 0")
	}

	if mo.WaitInterval <= 0 {
		return errors.New("wait_interval must be greater than 0")
	}

	return nil
}

func (so *MutateOptions) BindFlags(flags *pflag.FlagSet) error {
	return nil
}
