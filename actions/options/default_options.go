package options

import (
	"errors"
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type emailNotification struct {
	From    string `json:"from"`
	Subject string `json:"subject"`
}

type slackNotification struct {
	Webhook       string   `json:"webhook"`
	Color         string   `json:"color"`
	Fallback      string   `json:"fallback"`
	AuthorName    string   `json:"author_name"`
	AuthorSubname string   `json:"author_subname"`
	AuthorIcon    string   `json:"author_icon"`
	Footer        string   `json:"footer"`
	FooterIcon    string   `json:"footer_icon"`
	To            []string `json:"to"`
}

type notifiableOptions struct {
	emailNotification
	slackNotification
	NotificationChannel string `json:"notification_channel"`
}

type defaultOptions struct {
	notifiableOptions
	container     *gabs.Container
	Timeout       int  `json:"timeout_override"`
	DisableAction bool `json:"disable_action"`
}

func (do *defaultOptions) GetContainer() *gabs.Container {
	return do.container
}

func (do *defaultOptions) Exists(value string) bool {
	return do.container.Exists(value)
}

func (do *defaultOptions) GetDisableAction() bool {
	return do.DisableAction
}

func (do *defaultOptions) TimeoutInSeconds() string {
	if do.Timeout > 0 {
		return fmt.Sprintf("%vs", do.Timeout)
	}

	return ""
}

func (do *defaultOptions) Get(value string) interface{} {
	return do.container.S(value).Data()
}

func (do *defaultOptions) String(value string) string {
	return fmt.Sprint(do.Get(value))
}

func (do *defaultOptions) IsSnapshot() bool {
	return false
}

func (do defaultOptions) ValidateNotifiableOptions() error {
	if do.NotificationChannel == "email" {
		if len(do.From) == 0 || len(do.To) == 0 || len(do.Subject) == 0 {
			return errors.New("from, to, subject are required when specifying the email channel")
		}
	}

	if do.NotificationChannel == "slack" && len(do.Webhook) == 0 {
		return errors.New("a webhook is required when specifying the slack channel")
	}

	return nil
}

func (do *defaultOptions) IsNotifiable() bool {
	return inStringSlice(do.NotificationChannel, []string{"slack", "email"})
}

func (do *defaultOptions) defaultBindFlags(flags *pflag.FlagSet) error {
	do.Timeout = intFromFlags(flags, "timeout")
	do.DisableAction = boolFromFlags(flags, "disable_action")

	return nil
}
