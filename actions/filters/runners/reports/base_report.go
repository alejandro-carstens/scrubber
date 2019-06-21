package reports

import (
	"encoding/json"
	"fmt"
	"scrubber/actions/criterias"
	"time"

	"github.com/Jeffail/gabs"
)

type baseReport struct {
	FilterType string                `json:"filter_type"`
	Type       string                `json:"type"`
	Criteria   criterias.Criteriable `json:"criteria"`
	Summary    []string              `json:"summary"`
}

func (br *baseReport) SetType(actionType string) {
	br.Type = actionType
}

func (br *baseReport) SetCriteria(criteria criterias.Criteriable) {
	br.Criteria = criteria
	br.FilterType = criteria.Name()
}

func (br *baseReport) ToJson() (*gabs.Container, error) {
	b, err := json.Marshal(br)

	if err != nil {
		return nil, err
	}

	return gabs.ParseJSON(b)
}

func (br *baseReport) AddReason(reason string, values ...interface{}) {
	if len(br.Summary) == 0 {
		br.Summary = []string{}
	}

	reason = fmt.Sprintf(reason, values...)

	br.Summary = append(br.Summary, fmt.Sprintf("\t- [%v] %v", time.Now().Format(time.RFC3339), reason))
}

func (br *baseReport) toJsonString(value interface{}) (string, error) {
	b, err := json.Marshal(value)

	if err != nil {
		return "", err
	}

	container, err := gabs.ParseJSON(b)

	if err != nil {
		return "", err
	}

	return container.String(), nil
}
