package messages

import (
	"encoding/json"
	"testing"

	"github.com/Jeffail/gabs"
	"github.com/stretchr/testify/assert"
)

type contextTemplate struct {
	Count int
}

func TestEmailMessage(t *testing.T) {
	b, err := json.Marshal(map[string]interface{}{
		"notification_channel": "email",
		"to":                   []string{"Alejandro", "Scrubber", "Golang"},
		"from":                 "Alejandro",
		"subject":              "Scrubber Email Message Test",
		"text":                 "This is very cool {{ .Count }}",
	})

	assert.Nil(t, err)

	container, err := gabs.ParseJSON(b)

	assert.Nil(t, err)

	msg, err := NewMessage(container, &contextTemplate{Count: 10}, "random_dedup_key")

	assert.Nil(t, err)

	message, valid := msg.(*Email)

	assert.True(t, valid)
	assert.Equal(t, "Alejandro", message.From)
	assert.Equal(t, []string{"Alejandro", "Scrubber", "Golang"}, message.To)
	assert.Equal(t, "Scrubber Email Message Test", message.Subject)
	assert.Equal(t, "This is very cool 10", message.Body)
	assert.Equal(t, "random_dedup_key", message.DedupKey)
}

func TestPagerDutyMessage(t *testing.T) {
	b, err := json.Marshal(map[string]interface{}{
		"notification_channel": "pager_duty",
		"source":               "github.com/alejandro-carstens/scrubber",
		"severity":             "info",
		"component":            "application",
		"group":                "application",
		"class":                "application",
		"text":                 "This is very cool {{ .Count }}",
	})

	assert.Nil(t, err)

	container, err := gabs.ParseJSON(b)

	assert.Nil(t, err)

	msg, err := NewMessage(container, &contextTemplate{Count: 10}, "random_dedup_key")

	assert.Nil(t, err)

	message, valid := msg.(*PagerDuty)

	assert.True(t, valid)
	assert.Equal(t, "trigger", message.Event.EventAction)
	assert.Equal(t, "random_dedup_key", message.Event.DedupKey)
	assert.Equal(t, "github.com/alejandro-carstens/scrubber", message.Event.Payload.Source)
	assert.Equal(t, "info", message.Event.Payload.Severity)
	assert.Equal(t, "application", message.Event.Payload.Component)
	assert.Equal(t, "application", message.Event.Payload.Group)
	assert.Equal(t, "application", message.Event.Payload.Class)
	assert.Equal(t, "This is very cool 10", message.Event.Payload.Summary)
}

func TestSlackMessage(t *testing.T) {
	b, err := json.Marshal(map[string]interface{}{
		"notification_channel": "slack",
		"webhook":              "https://slack.webhook.com/test",
		"color":                "green",
		"fallback":             "fallback_test",
		"author_name":          "Alejandro",
		"author_subname":       "Alex",
		"author_icon":          "not available",
		"footer":               "test_footer",
		"to":                   []string{"Alejandro", "Carstens", "Cattori"},
		"footer_icon":          "test_footer_icon",
		"text":                 "This is very cool {{ .Count }}",
	})

	assert.Nil(t, err)

	container, err := gabs.ParseJSON(b)

	assert.Nil(t, err)

	msg, err := NewMessage(container, &contextTemplate{Count: 10}, "random_dedup_key")

	assert.Nil(t, err)

	message, valid := msg.(*Slack)

	assert.True(t, valid)

	assert.Equal(t, "green", message.Attachment.Color)
	assert.Equal(t, "fallback_test", message.Attachment.Fallback)
	assert.Equal(t, "Alejandro", message.Attachment.AuthorName)
	assert.Equal(t, "Alex", message.Attachment.AuthorSubname)
	assert.Equal(t, "not available", message.Attachment.AuthorIcon)
	assert.Equal(t, "Alejandro, Carstens, Cattori: This is very cool 10", message.Attachment.Text)
	assert.Equal(t, "test_footer", message.Attachment.Footer)
	assert.Equal(t, "test_footer_icon", message.Attachment.FooterIcon)
	assert.Equal(t, "https://slack.webhook.com/test", message.Attachment.WebhookName)
	assert.Equal(t, []string{"Alejandro", "Carstens", "Cattori"}, message.Attachment.To)
}
