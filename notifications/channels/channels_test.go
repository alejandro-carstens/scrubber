package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"scrubber/notifications/configurations"
	"scrubber/notifications/messages"

	"github.com/Jeffail/gabs"
	"github.com/go-mail/mail"
	"github.com/stretchr/testify/assert"
)

var slackRetryCounter int = -1
var pagerDutyRetryCounter int = -1

type contextTemplate struct {
	Count int
}

type mockSender mail.SendFunc

func (s mockSender) Send(from string, to []string, msg io.WriterTo) error {
	return s(from, to, msg)
}

type mockSendCloser struct {
	mockSender
	close func() error
}

func (s *mockSendCloser) Close() error {
	return s.close()
}

func TestSlackChannel(t *testing.T) {
	ts := setHttpTestMockServer(t)

	defer ts.Close()

	msg, err := setSlackMessage()

	assert.Nil(t, err)

	os.Setenv("SLACK_WEBHOOKS", fmt.Sprintf("team@%v/slack_webhook/random_token", ts.URL))

	assert.Nil(t, Notify(msg))
}

func TestSlackChannelRetry(t *testing.T) {
	ts := setHttpTestMockServer(t)

	defer ts.Close()

	os.Setenv("SLACK_WEBHOOKS", fmt.Sprintf("team@%v/slack_webhook_error", ts.URL))
	os.Setenv("SLACK_RETRY_COUNT", "3")

	msg, err := setSlackMessage()

	assert.Nil(t, err)

	assert.NotNil(t, Notify(msg))
	assert.Equal(t, 3, slackRetryCounter)
}

func TestPagerDutyChannel(t *testing.T) {
	ts := setHttpTestMockServer(t)

	defer ts.Close()

	msg, err := setPagerDutyMessage()

	assert.Nil(t, err)

	os.Setenv("PAGER_DUTY_ROUTING_KEY", "pager_duty_test_routing_key")
	os.Setenv("PAGER_DUTY_BASE_URI", fmt.Sprintf("%v/v2/enqueue", ts.URL))

	assert.Nil(t, Notify(msg))
}

func TestPagerDutyChannelRetry(t *testing.T) {
	ts := setHttpTestMockServer(t)

	defer ts.Close()

	msg, err := setPagerDutyMessage()

	assert.Nil(t, err)

	os.Setenv("PAGER_DUTY_ROUTING_KEY", "pager_duty_test_routing_key")
	os.Setenv("PAGER_DUTY_BASE_URI", fmt.Sprintf("%v/v2/enqueue/error", ts.URL))
	os.Setenv("PAGER_DUTY_RETRY_COUNT", "3")

	assert.NotNil(t, Notify(msg))
	assert.Equal(t, 3, pagerDutyRetryCounter)
}

func TestEmailChannelWithRetry(t *testing.T) {
	message, err := setEmailMessage()
	assert.Nil(t, err)

	os.Setenv("EMAIL_RETRY_COUNT", "2")
	os.Setenv("SMTP_HOST", "smtp_test_host")
	os.Setenv("SMTP_PORT", "90")
	os.Setenv("EMAIL_USERNAME", "test_email_username")
	os.Setenv("EMAIL_PASSWORD", "test_password")

	sendCloser := &mockSendCloser{
		mockSender: stubSend(t, "alejandro@test.com", []string{"alejandro@test.com", "scrubber@test.com"}, "This is very cool 10"),
		close: func() error {
			return nil
		},
	}

	channel := &Email{}
	channel.setSendCloser(sendCloser)

	config, err := configurations.Config(message.Type())

	assert.Nil(t, err)
	assert.Nil(t, channel.Configure(config))
	assert.Nil(t, channel.Send(message))
	assert.Nil(t, channel.Retry())
}

func setSlackMessage() (messages.Sendable, error) {
	b, err := json.Marshal(map[string]interface{}{
		"notification_channel": "slack",
		"webhook":              "team",
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

	if err != nil {
		return nil, err
	}

	container, err := gabs.ParseJSON(b)

	if err != nil {
		return nil, err
	}

	return messages.NewMessage(container, &contextTemplate{Count: 10}, "random_dedup_key")
}

func setPagerDutyMessage() (messages.Sendable, error) {
	b, err := json.Marshal(map[string]interface{}{
		"notification_channel": "pager_duty",
		"source":               "scrubber",
		"severity":             "info",
		"component":            "application",
		"group":                "application",
		"class":                "application",
		"text":                 "This is very cool {{ .Count }}",
	})

	if err != nil {
		return nil, err
	}

	container, err := gabs.ParseJSON(b)

	if err != nil {
		return nil, err
	}

	return messages.NewMessage(container, &contextTemplate{Count: 10}, "random_dedup_key")
}

func setEmailMessage() (messages.Sendable, error) {
	b, err := json.Marshal(map[string]interface{}{
		"notification_channel": "email",
		"to":                   []string{"alejandro@test.com", "scrubber@test.com"},
		"from":                 "alejandro@test.com",
		"subject":              "Scrubber Email Message Test",
		"text":                 "This is very cool {{ .Count }}",
	})

	if err != nil {
		return nil, err
	}

	container, err := gabs.ParseJSON(b)

	if err != nil {
		return nil, err
	}

	return messages.NewMessage(container, &contextTemplate{Count: 10}, "random_dedup_key")
}

func setHttpTestMockServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/slack_webhook/random_token" {
			request, err := parseRequestToJSON(r)

			if err != nil {
				t.Error(err)
			}

			attachments, err := request.S("attachments").Children()

			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, 1, len(attachments))
			assert.Equal(t, "green", attachments[0].S("color").Data().(string))
			assert.Equal(t, "Alejandro", attachments[0].S("author_name").Data().(string))
			assert.Equal(t, "Alejandro, Carstens, Cattori: This is very cool 10", attachments[0].S("text").Data().(string))

			w.WriteHeader(200)
		}

		if r.URL.Path == "/slack_webhook_error" {
			slackRetryCounter++

			w.WriteHeader(500)
		}

		if r.URL.Path == "/v2/enqueue" {
			request, err := parseRequestToJSON(r)

			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, "random_dedup_key", request.S("dedup_key").Data().(string))
			assert.Equal(t, "trigger", request.S("event_action").Data().(string))
			assert.Equal(t, "scrubber", request.S("payload", "source").Data().(string))
			assert.Equal(t, "info", request.S("payload", "severity").Data().(string))
			assert.Equal(t, "This is very cool 10", request.S("payload", "summary").Data().(string))
			assert.Equal(t, "application", request.S("payload", "component").Data().(string))

			w.WriteHeader(200)
		}

		if r.URL.Path == "/v2/enqueue/error" {
			pagerDutyRetryCounter++

			w.WriteHeader(500)
		}
	}))
}

func parseRequestToJSON(request *http.Request) (*gabs.Container, error) {
	b, err := ioutil.ReadAll(request.Body)

	defer request.Body.Close()

	if err != nil {
		return nil, err
	}

	return gabs.ParseJSON(b)
}

func stubSend(t *testing.T, wantFrom string, wantTo []string, wantBody string) mockSender {
	return func(from string, to []string, msg io.WriterTo) error {
		assert.Equal(t, from, wantFrom)
		assert.True(t, reflect.DeepEqual(to, wantTo))

		buf := new(bytes.Buffer)

		if _, err := msg.WriteTo(buf); err != nil {
			t.Fatal(err)
		}

		assert.True(t, strings.Contains(buf.String(), wantBody))

		return nil
	}
}
