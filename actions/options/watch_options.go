package options

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/spf13/pflag"
)

type slackAlert struct {
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

type Alert struct {
	slackAlert
	NotificationChannel string `json:"notification_channel"`
	Text                string `json:"text"`
}

func (a *Alert) validate() error {
	if len(a.NotificationChannel) == 0 {
		return errors.New("an alert requires a notification_channel")
	}

	if len(a.Text) == 0 {
		return errors.New("an alert requires some text")
	}

	if a.NotificationChannel == "slack" && len(a.Webhook) == 0 {
		return errors.New("a webhook is required when specifying the slack channel")
	}

	return nil
}

func (a *Alert) Payload() *gabs.Container {
	return toContainer(a)
}

type Threshold struct {
	Type   string   `json:"type"`
	Metric string   `json:"metric"`
	Min    *float64 `json:"min"`
	Max    *float64 `json:"max"`
	Alerts []*Alert `json:"alerts"`
}

func (t *Threshold) validate(dateField, statsField string) error {
	if t.Max == nil && t.Min == nil {
		return errors.New("a min limit or a max limit need to be specified")
	}

	if !inStringSlice(t.Type, availableThresholdTypes) {
		return errors.New("invalid threshold type specified")
	}

	if t.Type == "stats" && !inStringSlice(t.Metric, availableMetrics) {
		return errors.New("specified metric is not compatible with the stats threshold type")
	}

	if t.Type == "average_count" && len(dateField) == 0 {
		return errors.New("date_field is required for the average_count threshold")
	}

	if t.Type == "stats" && len(statsField) == 0 {
		return errors.New("stats_field is required when specifying a stats threshold type")
	}

	for _, alert := range t.Alerts {
		if err := alert.validate(); err != nil {
			return err
		}
	}

	return nil
}

type queryCriteria struct {
	Clause   string        `json:"clause"`
	Key      string        `json:"key"`
	Operator string        `json:"operator"`
	Value    interface{}   `json:"value"`
	Values   []interface{} `json:"values"`
	Limit    int           `json:"limit"`
	Order    bool          `json:"order"`
}

func (qc *queryCriteria) validate() error {
	if len(qc.Key) == 0 {
		return errors.New("key on query criteria cannot be empty")
	}

	if !inStringSlice(qc.Clause, availableClauses) {
		return errors.New("invalida clause specified")
	}

	if len(qc.Operator) > 0 && !inStringSlice(qc.Operator, availableOperators) {
		return errors.New("invalid operator specified")
	}

	if inStringSlice(qc.Clause, availableInClauses) && len(qc.Values) == 0 {
		return errors.New("values param is required when using 'In' cluases")
	}

	if !inStringSlice(qc.Clause, append(availableInClauses, nonMatchingClauses...)) && qc.Value == nil {
		return errors.New("value is a required param")
	}

	if qc.Clause == "limit" && qc.Limit <= 0 {
		return errors.New("limit param needs to be greater than 0")
	}

	return nil
}

type WatchOptions struct {
	defaultOptions
	Interval        int64            `json:"interval"`
	IntervalUnit    string           `json:"interval_unit"`
	DateField       string           `json:"date_field"`
	StatsField      string           `json:"stats_field"`
	Criteria        []*queryCriteria `json:"criteria"`
	Thresholds      []*Threshold     `json:"thresholds"`
	AlertChannels   []string         `json:"alert_channels"`
	MessageTemplate string           `json:"message_template"`
}

func (wo *WatchOptions) FillFromContainer(container *gabs.Container) error {
	wo.container = container

	return json.Unmarshal(container.Bytes(), wo)
}

func (wo *WatchOptions) Validate() error {
	for _, criteria := range wo.Criteria {
		if err := criteria.validate(); err != nil {
			return err
		}
	}

	for _, threshold := range wo.Thresholds {
		if err := threshold.validate(wo.DateField, wo.StatsField); err != nil {
			return err
		}
	}

	for _, alertChannel := range wo.AlertChannels {
		if !inStringSlice(alertChannel, availableAlertChannels) {
			return errors.New("invalid alert channel")
		}
	}

	if !inStringSlice(wo.IntervalUnit, availableIntervalUnits) {
		return errors.New("ivalid interval unit")
	}

	return nil
}

func (so *WatchOptions) BindFlags(flags *pflag.FlagSet) error {
	return nil
}
