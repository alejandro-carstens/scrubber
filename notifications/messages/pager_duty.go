package messages

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"

	"github.com/Jeffail/gabs"
)

type eventPayload struct {
	Source    string `json:"source"`
	Severity  string `json:"severity"`
	Component string `json:"component"`
	Group     string `json:"group"`
	Class     string `json:"class"`
	Summary   string `json:"summary"`
}

type pagerDutyEvent struct {
	DedupKey    string        `json:"dedup_key"`
	EventAction string        `json:"event_action"`
	Payload     *eventPayload `json:"payload"`
}

// Pager Duty is the representation of a Pager Duty message
type PagerDuty struct {
	Context  interface{}
	Payload  *gabs.Container
	Event    *pagerDutyEvent
	DedupKey string
}

// Type returns the message type
func (pd *PagerDuty) Type() string {
	return "pager_duty"
}

// Format formats a message to be sent over the specified channel
func (pd *PagerDuty) Format() error {
	summary, valid := pd.Payload.S("text").Data().(string)

	if !valid {
		return errors.New("text field not found in payload")
	}

	eventPayload := &eventPayload{Summary: summary}

	if err := json.Unmarshal(pd.Payload.Bytes(), eventPayload); err != nil {
		return err
	}

	event := &pagerDutyEvent{DedupKey: pd.DedupKey, EventAction: "trigger", Payload: eventPayload}

	templ, err := template.New("text").Parse(event.Payload.Summary)

	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)

	if err := templ.Execute(buffer, pd.Context); err != nil {
		return err
	}

	event.Payload.Summary = buffer.String()

	pd.Event = event

	return nil
}
