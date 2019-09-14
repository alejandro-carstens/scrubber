package actions

import (
	"encoding/json"
	"errors"
	"scrubber/actions/options"
	"scrubber/notifications"
	"scrubber/notifications/messages"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/golastic"
)

type count struct {
	isSet bool
	Count float64 `json:"count"`
}

type stats struct {
	isSet                  bool
	Min                    float64 `json:"min"`
	Max                    float64 `json:"max"`
	Avg                    float64 `json:"avg"`
	Sum                    float64 `json:"sum"`
	SumOfSquares           float64 `json:"sum_of_squares"`
	Variance               float64 `json:"variance"`
	StdDeviation           float64 `json:"std_deviation"`
	UpperStdDeviationBound float64 `json:"upper_std_deviation_bound"`
	LowerStdDeviationBound float64 `json:"lower_std_deviation_bound"`
}

type watch struct {
	filterAction
	options *options.WatchOptions
}

func (w *watch) ApplyOptions() Actionable {
	w.options = w.context.Options().(*options.WatchOptions)

	w.indexer.SetOptions(&golastic.IndexOptions{Timeout: w.options.TimeoutInSeconds()})

	return w
}

func (w *watch) Perform() Actionable {
	w.exec(func(index string) error {
		var err error
		var mappings *gabs.Container

		if len(w.options.DateField)+len(w.options.StatsField) > 1 {
			mappings, err = w.indexer.FieldMappings(index)

			if err != nil {
				return err
			}
		}

		if len(w.options.DateField) > 0 {
			path := []string{index, "mappings", w.options.DateField, "mapping", w.options.DateField, "type"}

			mappingType, valid := mappings.S(path...).Data().(string)

			if !valid || mappingType != "date" {
				return errors.New("invalid date_field specified")
			}
		}

		if len(w.options.StatsField) > 0 {
			path := []string{index, "mappings", w.options.StatsField, "mapping", w.options.StatsField, "type"}

			mappingType, valid := mappings.S(path...).Data().(string)

			if !valid || !inStringSlice(mappingType, availableNumericTypes) {
				return errors.New("invalid stats_field specified")
			}
		}

		if err := w.execute(index); err != nil {
			return err
		}

		return nil
	})

	return w
}

func (w *watch) execute(index string) error {
	count := &count{}
	stats := &stats{}

	builder := w.buildQuery(index)

	for _, threshold := range w.options.Thresholds {
		if inStringSlice(threshold.Type, []string{"count", "average_count"}) && !count.isSet {
			value, err := builder.Count()

			if err != nil {
				return err
			}

			count.Count = float64(value)
			count.isSet = true
		} else if threshold.Type == "stats" && !stats.isSet {
			response, err := builder.AggregateRaw()

			if err != nil {
				return err
			}

			if err := json.Unmarshal(response.S(w.options.StatsField).Bytes(), stats); err != nil {
				return err
			}

			stats.isSet = true
		}

		var err error

		switch threshold.Type {
		case "count":
			err = w.processCountThreshold(count, threshold)
			break
		case "average_count":
			err = w.processAverageCountThreshold(count, threshold)
			break
		case "stats":
			err = w.processStats(stats, threshold)
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (w *watch) buildQuery(index string) *golastic.Builder {
	builder := w.connection.Builder(index)

	for _, criteria := range w.options.Criteria {
		switch criteria.Clause {
		case "where":
			builder.Where(criteria.Key, criteria.Operator, criteria.Value)
			break
		case "where_nested":
			builder.WhereNested(criteria.Key, criteria.Operator, criteria.Value)
			break
		case "where_in":
			builder.WhereIn(criteria.Key, criteria.Values)
			break
		case "where_in_nested":
			builder.WhereInNested(criteria.Key, criteria.Values)
			break
		case "where_not_in":
			builder.WhereNotIn(criteria.Key, criteria.Values)
			break
		case "where_not_in_nested":
			builder.WhereNotInNested(criteria.Key, criteria.Values)
			break
		case "filter":
			builder.Filter(criteria.Key, criteria.Operator, criteria.Value)
			break
		case "filter_nested":
			builder.FilterNested(criteria.Key, criteria.Operator, criteria.Value)
			break
		case "filter_in":
			builder.FilterIn(criteria.Key, criteria.Values)
			break
		case "filter_in_nested":
			builder.FilterInNested(criteria.Key, criteria.Values)
			break
		case "match":
			builder.Match(criteria.Key, criteria.Operator, criteria.Value)
			break
		case "match_nested":
			builder.MatchNested(criteria.Key, criteria.Operator, criteria.Value)
			break
		case "match_in":
			builder.MatchIn(criteria.Key, criteria.Values)
			break
		case "match_in_nested":
			builder.MatchInNested(criteria.Key, criteria.Values)
			break
		case "match_not_in":
			builder.MatchNotIn(criteria.Key, criteria.Values)
			break
		case "match_not_in_nested":
			builder.MatchNotInNested(criteria.Key, criteria.Values)
			break
		case "limit":
			builder.Limit(criteria.Limit)
			break
		case "order_by":
			builder.OrderBy(criteria.Key, criteria.Order)
			break
		case "order_by_nested":
			builder.OrderByNested(criteria.Key, criteria.Order)
			break
		}
	}

	if len(w.options.DateField) > 0 {
		duration := -1 * intervalToSeconds(w.options.Interval, w.options.IntervalUnit)

		builder.Filter(w.options.DateField, ">=", time.Now().Add(time.Duration(duration)*time.Second))
	}

	if len(w.options.StatsField) > 0 {
		builder.Stats(w.options.StatsField)
	}

	return builder
}

func (w *watch) processCountThreshold(count *count, threshold *options.Threshold) error {
	return w.compare(float64(count.Count), threshold, count)
}

func (w *watch) processAverageCountThreshold(count *count, threshold *options.Threshold) error {
	count.Count = float64(count.Count) / float64(intervalToSeconds(w.options.Interval, w.options.IntervalUnit))

	return w.compare(count.Count, threshold, count)
}

func (w *watch) processStats(stats *stats, threshold *options.Threshold) error {
	switch threshold.Metric {
	case "min":
		return w.compare(stats.Min, threshold, stats)
	case "max":
		return w.compare(stats.Max, threshold, stats)
	case "avg":
		return w.compare(stats.Avg, threshold, stats)
	case "sum":
		return w.compare(stats.Sum, threshold, stats)
	case "sum_of_squares":
		return w.compare(stats.SumOfSquares, threshold, stats)
	case "variance":
		return w.compare(stats.Variance, threshold, stats)
	case "std_deviation":
		return w.compare(stats.StdDeviation, threshold, stats)
	case "upper_std_deviation_bound":
		return w.compare(stats.UpperStdDeviationBound, threshold, stats)
	case "lower_std_deviation_bound":
		return w.compare(stats.LowerStdDeviationBound, threshold, stats)
	}

	return nil
}

func (w *watch) compare(metric float64, threshold *options.Threshold, context interface{}) error {
	min := *threshold.Min
	max := *threshold.Max

	w.reporter.Logger().Noticef("metric: %v, min: %v, max: %v", metric, min, max)

	if threshold.Min != nil && metric < min {
		return w.alert(threshold.Alerts, context)
	}

	if threshold.Max != nil && metric > max {
		return w.alert(threshold.Alerts, context)
	}

	return nil
}

func (w *watch) alert(alerts []*options.Alert, context interface{}) error {
	for _, alert := range alerts {
		var err error

		message, err := messages.NewMessage(alert.Payload(), context)

		if err != nil {
			return err
		}

		if w.enqueuer != nil {
			err = w.enqueuer.Push(message)
		} else {
			err = notifications.Notify(message)
		}

		if err != nil {
			return err
		}
	}

	return nil
}