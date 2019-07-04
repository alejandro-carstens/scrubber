package criterias

import (
	"errors"

	"github.com/Jeffail/gabs"
)

type Builder struct {
	criteria          []Criteriable
	aggregateCriteria []Criteriable
}

func (b *Builder) Build(action string, filters []*gabs.Container) error {
	criteria := []Criteriable{}
	aggregateCriteria := []Criteriable{}

	for _, filter := range filters {
		filterType, valid := filter.S("filtertype").Data().(string)

		if !valid {
			return errors.New("No filtertype specified for filter")
		}

		if b.isSnapshotAction(action) && !b.isSnapshotFilterType(filterType) {
			return errors.New("Invalid filter " + filterType + " for snapshot action")
		}

		filterCriteria, err := b.fillCriteria(filterType, filter)

		if err != nil {
			return err
		}

		if filterCriteria != nil {
			criteria = append(criteria, filterCriteria)

			continue
		}

		filterCriteria, err = b.fillAggregateCriteria(filterType, filter)

		if err != nil {
			return err
		}

		aggregateCriteria = append(aggregateCriteria, filterCriteria)
	}

	b.criteria = criteria
	b.aggregateCriteria = aggregateCriteria

	return nil
}

func (b *Builder) Criteria() []Criteriable {
	return b.criteria
}

func (b *Builder) AggregateCriteria() []Criteriable {
	return b.aggregateCriteria
}

func (b *Builder) fillCriteria(filterType string, filter *gabs.Container) (Criteriable, error) {
	var criteria Criteriable
	var err error

	switch filterType {
	case "age":
		criteria, err = new(Age).FillFromContainer(filter)
		break
	case "pattern":
		criteria, err = new(Pattern).FillFromContainer(filter)
		break
	case "empty":
		criteria, err = new(Empty).FillFromContainer(filter)
		break
	case "kibana":
		criteria, err = new(Kibana).FillFromContainer(filter)
		break
	case "closed":
		criteria, err = new(Closed).FillFromContainer(filter)
		break
	case "alias":
		criteria, err = new(Alias).FillFromContainer(filter)
		break
	case "allocated":
		criteria, err = new(Allocated).FillFromContainer(filter)
		break
	case "forcemerged":
		criteria, err = new(Forcemerged).FillFromContainer(filter)
		break
	case "state":
		criteria, err = new(State).FillFromContainer(filter)
		break
	}

	if err != nil || criteria == nil {
		return nil, err
	}

	return criteria, criteria.Validate()
}

func (b *Builder) fillAggregateCriteria(filterType string, filter *gabs.Container) (Criteriable, error) {
	var criteria Criteriable
	var err error

	switch filterType {
	case "count":
		criteria, err = new(Count).FillFromContainer(filter)
		break
	case "space":
		criteria, err = new(Space).FillFromContainer(filter)
		break
	default:
		return nil, errors.New("No valid filtertype specified")
	}

	if err != nil {
		return nil, err
	}

	return criteria, criteria.Validate()
}

func (b *Builder) isSnapshotAction(action string) bool {
	switch action {
	case "restore":
		return true
	case "delete_snapshot":
		return true
	case "list_snapshots":
		return true
	}

	return false
}

func (b *Builder) isSnapshotFilterType(filterType string) bool {
	switch filterType {
	case "age":
		return true
	case "pattern":
		return true
	case "count":
		return true
	case "state":
		return true
	}

	return false
}
