package actions

import (
	"errors"
	"scrubber/actions/options"
	"time"

	"github.com/alejandro-carstens/golastic"
)

type watch struct {
	filterAction
	options *options.WatchOptions
}

type count struct {
	count int64
	isSet bool
}

func (w *watch) ApplyOptions() Actionable {
	w.options = w.context.Options().(*options.WatchOptions)

	w.indexer.SetOptions(&golastic.IndexOptions{Timeout: w.options.TimeoutInSeconds()})

	return w
}

func (w *watch) Perform() Actionable {
	w.exec(func(index string) error {
		if len(w.options.DateField) > 0 {
			mappings, err := w.indexer.FieldMappings(index)

			if err != nil {
				return err
			}

			path := []string{index, "mappings", w.options.DateField, "mapping", w.options.DateField, "type"}

			mappingType, valid := mappings.S(path...).Data().(string)

			if !valid || mappingType != "date" {
				return errors.New("invalid date_field specified")
			}
		}

		builder := w.buildQuery(index)

		if err := w.execute(builder); err != nil {
			return err
		}

		return nil
	})

	return w
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

		builder.Where(w.options.DateField, ">=", time.Now().Add(time.Duration(duration)*time.Second))
	}

	return builder
}

func (w *watch) execute(builder *golastic.Builder) error {
	count := &count{}

	for _, threshold := range w.options.Thresholds {
		if threshold.Type == "count" || threshold.Type == "average_count" {
			if !count.isSet {
				value, err := builder.Count()

				if err != nil {
					return err
				}

				count.count = value
				count.isSet = true
			}
		}

		var err error

		switch threshold.Type {
		case "count":
			err = w.processCountThreshold(count.count, threshold)
			break
		case "average_count":
			err = w.processAverageCountThreshold(count.count, threshold)
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (w *watch) processCountThreshold(count int64, threshold *options.Threshold) error {
	if float64(count) < threshold.Min {
		w.reporter.Logger().Noticef("min threshold of %v exceeded, encountered %v", threshold.Min, count)

		return nil
	}

	if float64(count) > threshold.Max {
		w.reporter.Logger().Noticef("max threshold of %v exceeded, encountered %v", threshold.Max, count)

		return nil
	}

	return nil
}

func (w *watch) processAverageCountThreshold(count int64, threshold *options.Threshold) error {
	averageCount := float64(count) / float64(intervalToSeconds(w.options.Interval, w.options.IntervalUnit))

	if averageCount < threshold.Min {
		w.reporter.Logger().Noticef("min threshold of %v exceeded, encountered %v", threshold.Min, averageCount)

		return nil
	}

	if averageCount > threshold.Max {
		w.reporter.Logger().Noticef("max threshold of %v exceeded, encountered %v", threshold.Max, averageCount)

		return nil
	}

	return nil
}
