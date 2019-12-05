package messages

import (
	"encoding/json"
	"testing"

	"github.com/Jeffail/gabs"
	"github.com/stretchr/testify/assert"
)

type ContextTemplate struct {
	Count int
}

func TestEmailMessage(t *testing.T) {
	b, err := json.Marshal(map[string]interface{}{
		"to":      []string{"Alejandro", "Scrubber", "Golang"},
		"from":    "Alejandro",
		"subject": "Scrubber Email Message Test",
		"text":    "This is very cool {{ .Count }}",
	})

	if err != nil {
		t.Error(err)
	}

	container, err := gabs.ParseJSON(b)

	if err != nil {
		t.Error(err)
	}

	message := &Email{}
	message.Init(&ContextTemplate{Count: 10}, "random_dedup_key", container)

	assert.Nil(t, message.Format())
	assert.Equal(t, "Alejandro", message.From)
	assert.Equal(t, []string{"Alejandro", "Scrubber", "Golang"}, message.To)
	assert.Equal(t, "Scrubber Email Message Test", message.Subject)
	assert.Equal(t, "This is very cool 10", message.Body)
	assert.Equal(t, "random_dedup_key", message.DedupKey)
}

func TestPagerDutyMessage(t *testing.T) {
	b, err := json.Marshal(map[string]interface{}{
		"source":    "scrubber",
		"severity":  "info",
		"component": "application",
		"group":     "application",
		"class":     "application",
		"text":      "This is very cool {{ .Count }}",
	})

	if err != nil {
		t.Error(err)
	}

	container, err := gabs.ParseJSON(b)

	if err != nil {
		t.Error(err)
	}

	message := &PagerDuty{}
	message.Init(&ContextTemplate{Count: 10}, "random_dedup_key", container)

	assert.Nil(t, message.Format())
	assert.Equal(t, "trigger", message.Event.EventAction)
	assert.Equal(t, "random_dedup_key", message.Event.DedupKey)
	assert.Equal(t, "scrubber", message.Event.Payload.Source)
	assert.Equal(t, "info", message.Event.Payload.Severity)
	assert.Equal(t, "application", message.Event.Payload.Component)
	assert.Equal(t, "application", message.Event.Payload.Group)
	assert.Equal(t, "application", message.Event.Payload.Class)
	assert.Equal(t, "This is very cool 10", message.Event.Payload.Summary)
}

func TestSlackMessage(t *testing.T) {
	b, err := json.Marshal(map[string]interface{}{
		"webhook":        "https://slack.webhook.com/test",
		"color":          "green",
		"fallback":       "fallback_test",
		"author_name":    "Alejandro",
		"author_subname": "Alex",
		"author_icon":    "not available",
		"footer":         "test_footer",
		"footer_icon":    "test_footer_icon",
		"to":             []string{"Alejandro", "Carstens", "Cattori"},
		"text":           "This is very cool {{ .Count }}",
	})

	if err != nil {
		t.Error(err)
	}

	container, err := gabs.ParseJSON(b)

	if err != nil {
		t.Error(err)
	}

	message := &PagerDuty{}
	message.Init(&ContextTemplate{Count: 10}, "random_dedup_key", container)

	assert.Nil(t, message.Format())
}
