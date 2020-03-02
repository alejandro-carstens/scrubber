package actions

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
)

type actionNotification struct {
	NotificationChannel string   `json:"notification_channel"`
	Text                string   `json:"text"`
	Webhook             string   `json:"webhook"`
	Color               string   `json:"color"`
	Fallback            string   `json:"fallback"`
	AuthorName          string   `json:"author_name"`
	AuthorSubname       string   `json:"author_subname"`
	AuthorIcon          string   `json:"author_icon"`
	Footer              string   `json:"footer"`
	FooterIcon          string   `json:"footer_icon"`
	To                  []string `json:"to"`
	From                string   `json:"from"`
	Subject             string   `json:"subject"`
}

func (an *actionNotification) payload() *gabs.Container {
	b, _ := json.Marshal(an)

	container, _ := gabs.ParseJSON(b)

	return container
}
