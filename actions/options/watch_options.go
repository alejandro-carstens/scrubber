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

type queryCriteria struct {
	Clause   string        `json:"clause"`
	Key      string        `json:"key"`
	Operator string        `json:"operator"`
	Value    interface{}   `json:"value"`
	Values   []interface{} `json:"values"`
	Limit    int           `json:"limit"`
	Order    bool          `json:"order"`
}

type Threshold struct {
	Type   string   `json:"type"`
	Metric string   `json:"metric"`
	Min    *float64 `json:"min"`
	Max    *float64 `json:"max"`
	Alerts []*Alert `json:"alerts"`
}

type Alert struct {
	slackAlert
	NotificationChannel string `json:"notification_channel"`
	Text                string `json:"text"`
}

func (a *Alert) Payload() *gabs.Container {
	return toContainer(a)
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
	if err := wo.validateCriteria(); err != nil {
		return err
	}

	if err := wo.validateThresholds(); err != nil {
		return err
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

func (wo *WatchOptions) validateCriteria() error {
	for _, criteria := range wo.Criteria {
		if len(criteria.Key) == 0 {
			return errors.New("key on query criteria cannot be empty")
		}

		if !inStringSlice(criteria.Clause, availableClauses) {
			return errors.New("invalida clause specified")
		}

		if len(criteria.Operator) > 0 && !inStringSlice(criteria.Operator, availableOperators) {
			return errors.New("invalid operator specified")
		}

		if inStringSlice(criteria.Clause, availableInClauses) && len(criteria.Values) == 0 {
			return errors.New("values param is required when using 'In' cluases")
		}

		if !inStringSlice(criteria.Clause, append(availableInClauses, nonMatchingClauses...)) && criteria.Value == nil {
			return errors.New("value is a required param")
		}

		if criteria.Clause == "limit" && criteria.Limit <= 0 {
			return errors.New("limit param needs to be greater than 0")
		}
	}

	return nil
}

func (wo *WatchOptions) validateThresholds() error {
	for _, threshold := range wo.Thresholds {
		if threshold.Max == nil && threshold.Min == nil {
			return errors.New("a min limit or a max limit need to be specified")
		}

		if !inStringSlice(threshold.Type, availableThresholdTypes) {
			return errors.New("invalid threshold type specified")
		}

		if threshold.Type == "stats" && len(wo.StatsField) == 0 {
			return errors.New("stats_field is required when specifying a stats threshold type")
		}

		if threshold.Type == "stats" && !inStringSlice(threshold.Metric, availableMetrics) {
			return errors.New("specified metric is not compatible with the stats threshold type")
		}

		if threshold.Type == "average_count" && len(wo.DateField) == 0 {
			return errors.New("date_field is required for the average_count threshold")
		}

		for _, alert := range threshold.Alerts {
			if len(alert.NotificationChannel) == 0 {
				return errors.New("an alert requires a notification_channel")
			}

			if len(alert.Text) == 0 {
				return errors.New("an alert requires some text")
			}
		}
	}

	return nil
}
