package messages

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
)

// Email is the representation of an e-mail message
type Email struct {
	baseMessage
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

// Type returns the message type
func (e *Email) Type() string {
	return "email"
}

// Format sets the format of a Sendable message
func (e *Email) Format() error {
	if err := json.Unmarshal(e.Payload.Bytes(), e); err != nil {
		return err
	}

	body, valid := e.Payload.S("text").Data().(string)

	if !valid {
		return errors.New("text field not found in payload")
	}

	templ, err := template.New("text").Parse(body)

	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)

	if err := templ.Execute(buffer, e.Context); err != nil {
		return err
	}

	e.Body = buffer.String()

	return nil
}
