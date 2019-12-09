package channels

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"scrubber/notifications/messages"
	"testing"

	"github.com/Jeffail/gabs"
	"github.com/stretchr/testify/assert"
)

var slackRetryCounter int = -1

type contextTemplate struct {
	Count int
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
