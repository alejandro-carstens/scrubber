package messages

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"strings"

	"github.com/Jeffail/gabs"
)

type slackAttachment struct {
	Color         string   `json:"color"`
	Fallback      string   `json:"fallback"`
	AuthorName    string   `json:"author_name"`
	AuthorSubname string   `json:"author_subname"`
	AuthorIcon    string   `json:"author_icon"`
	Text          string   `json:"text"`
	Footer        string   `json:"footer"`
	FooterIcon    string   `json:"footer_icon"`
	WebhookName   string   `json:"webhook"`
	To            []string `json:"to"`
}

// Slack is the representation of a Slack message
type Slack struct {
	Context    interface{}
	Payload    *gabs.Container
	Attachment *slackAttachment
}

// Type returns the message type
func (s *Slack) Type() string {
	return "slack"
}

// Format formats a message to be sent over the specified channel
func (s *Slack) Format() error {
	attachment := &slackAttachment{}

	if err := json.Unmarshal(s.Payload.Bytes(), attachment); err != nil {
		return err
	}

	if len(attachment.Text) == 0 {
		return errors.New("the slack message requires a text")
	}

	if len(attachment.AuthorName) == 0 {
		attachment.AuthorName = "scrubber"
	}

	templ, err := template.New("text").Parse(attachment.Text)

	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)

	if err := templ.Execute(buffer, s.Context); err != nil {
		return err
	}

	text := strings.Join(attachment.To, ", ")

	if len(text) > 0 {
		text = text + " "
	}

	attachment.Text = text + buffer.String()

	s.Attachment = attachment

	return nil
}
