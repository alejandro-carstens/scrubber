package runners

import (
	"errors"
	"math"
	"regexp"
	"scrubber/actions/criterias"
	"scrubber/actions/infos"
	"strconv"
	"time"

	"github.com/alejandro-carstens/golastic"
	"github.com/araddon/dateparse"
)

type ageFilterRunner struct {
	baseRunner
	builder golastic.Queryable
}

// Init initializes the filter runner
func (afr *ageFilterRunner) Init(info ...infos.Informable) (Runnerable, error) {
	if err := afr.BaseInit(info...); err != nil {
		return nil, err
	}

	var builder *golastic.ElasticsearchBuilder
	var err error

	if !afr.info.IsSnapshotInfo() {
		model := golastic.NewGolasticModel()
		model.SetIndex(afr.info.Name())

		builder, err = golastic.NewBuilder(model, nil)
	} else {
		builder, err = golastic.NewBuilder(nil, nil)
	}

	if err != nil {
		return nil, err
	}

	afr.builder = builder

	return afr, nil
}

// RunFilter filters out elements from the actionable list
func (afr *ageFilterRunner) RunFilter(channel chan *FilterResponse, criteria criterias.Criteriable) {
	if err := afr.validateCriteria(criteria); err != nil {
		channel <- afr.response.setError(err)
		return
	}

	var passed bool
	var err error

	age := criteria.(*criterias.Age)

	switch age.Source {
	case "creation_date":
		passed, err = afr.processByCreationDate(age)
		break
	case "name":
		passed, err = afr.processByName(age)
		break
	case "field_stats":
		passed, err = afr.processByFieldStats(age)
		break
	}

	channel <- afr.response.setError(err).setPassed(passed && age.Include()).setReport(afr.report)
}

func (afr *ageFilterRunner) processByCreationDate(age *criterias.Age) (bool, error) {
	creationDate, err := afr.creationDate()

	if err != nil {
		return false, err
	}

	creationTimestamp, err := time.Parse(time.RFC3339, creationDate)

	if err != nil {
		return false, err
	}

	afr.report.AddReason("Filtering by creation date")
	afr.report.AddReason("Creation date: %v", creationTimestamp.Format(time.RFC3339))

	return afr.compare(creationTimestamp, age)
}

func (afr *ageFilterRunner) processByName(age *criterias.Age) (bool, error) {
	var date time.Time
	var err error

	switch age.Timestring {
	case "Y.m.d":
		date, err = afr.parseDateFromName(`\d{4}.\d{2}.\d{2}`)
		break
	case "m.d.Y":
		date, err = afr.parseDateFromName(`\d{2}.\d{2}.\d{4}`)
		break
	case "Y.m":
		date, err = afr.parseDateFromName(`\d{4}.\d{2}`)
		break
	case "Y-m-d":
		date, err = afr.parseDateFromName(`\d{4}-\d{2}-\d{2}`)
		break
	case "Y-m-d H:M":
		date, err = afr.parseDateFromName(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}`)
		break
	case "Y-m-d H:M:S":
		date, err = afr.parseDateFromName(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`)
		break
	case "m-d-Y":
		date, err = afr.parseDateFromName(`\d{2}-\d{2}-\d{4}`)
		break
	case "Y-m":
		date, err = afr.parseDateFromName(`\d{4}-\d{2}`)
		break
	}

	if err != nil {
		return false, err
	}

	afr.report.AddReason("Filtering by name")
	afr.report.AddReason(
		"Generated timestamp for index '%v' and timestring '%v': %v",
		afr.info.Name(),
		age.Timestring,
		date.Format(time.RFC3339),
	)

	return afr.compare(date, age)
}

func (afr *ageFilterRunner) processByFieldStats(age *criterias.Age) (bool, error) {
	if afr.info.IsSnapshotInfo() {
		return false, errors.New("Cannot process age filter by fields_stats for snapshot action")
	}

	result, err := afr.builder.MinMax(age.Field, true)

	if err != nil {
		return false, err
	}

	var date string
	var valid bool

	switch age.StatsResult {
	case "min":
		date, valid = result.Min.(string)
		break
	case "max":
		date, valid = result.Max.(string)
		break
	}

	if !valid {
		return false, errors.New("Field is not a valid date parameter")
	}

	dateTime, err := dateparse.ParseLocal(date)

	if err != nil {
		return false, err
	}

	afr.report.AddReason("Filtering by field stats")
	afr.report.AddReason(
		"Stats result '%v' for field '%v': %v",
		age.StatsResult,
		age.Field,
		dateTime.Format(time.RFC3339),
	)

	return afr.compare(dateTime, age)
}

func (afr *ageFilterRunner) compare(date time.Time, age *criterias.Age) (bool, error) {
	duration := -1 * int64(age.UnitCount)
	var since time.Time

	switch age.Units {
	case "seconds":
		since = time.Now().UTC().Add(time.Duration(duration) * time.Second)
		break
	case "minutes":
		since = time.Now().UTC().Add(time.Duration(duration) * time.Minute)
		break
	case "hours":
		since = time.Now().UTC().Add(time.Duration(duration) * time.Hour)
		break
	case "days":
		since = time.Now().UTC().AddDate(0, 0, int(duration))
		break
	case "months":
		since = time.Now().UTC().AddDate(int(duration), 0, 0)
		break
	case "years":
		since = time.Now().UTC().AddDate(int(duration), 0, 0)
		break
	}

	diff := afr.diff(since, date, age)

	afr.report.AddReason(
		"Comparison date based on the current time minus '%v' '%v': %v",
		strconv.Itoa(age.UnitCount),
		age.Units,
		since.Format(time.RFC3339),
	)

	switch age.Direction {
	case "older":
		return date.UTC().Before(since) && diff, nil
	}

	return date.UTC().After(since) && diff, nil
}

func (afr *ageFilterRunner) parseDateFromName(regPattern string) (time.Time, error) {
	reg, err := regexp.Compile(regPattern)

	if err != nil {
		return time.Now(), nil
	}

	value := reg.FindString(afr.info.Name())

	if len(value) == 0 {
		return time.Now(), errors.New("Could not match the expected pattern")
	}

	return dateparse.ParseLocal(value)
}

func (afr *ageFilterRunner) diff(from, to time.Time, age *criterias.Age) bool {
	years, months, days, hours, minutes, seconds := elapsed(from, to)

	afr.report.AddReason(
		"Diff in Years: %v, Months: %v, Days: %v, Hours: %v, Minutes: %v, Seconds: %v",
		years,
		months,
		days,
		hours,
		minutes,
		seconds,
	)

	switch age.Units {
	case "seconds":
		return math.Abs(from.UTC().Sub(to.UTC()).Seconds()) > float64(age.UnitCount)
	case "minutes":
		return math.Abs(from.UTC().Sub(to.UTC()).Minutes()) > float64(age.UnitCount)
	case "hours":
		return math.Abs(from.UTC().Sub(to.UTC()).Hours()) > float64(age.UnitCount)
	case "days":
		return secondsToDays(math.Abs(from.UTC().Sub(to.UTC()).Seconds())) > float64(age.UnitCount)
	case "months":
		return secondsToMonths(math.Abs(from.UTC().Sub(to.UTC()).Seconds())) > float64(age.UnitCount)
	}

	return secondsToYears(math.Abs(from.UTC().Sub(to.UTC()).Seconds())) > float64(age.UnitCount)
}

func (afr *ageFilterRunner) creationDate() (string, error) {
	creationDate := afr.info.CreationDate()

	if len(creationDate) == 0 {
		return "", errors.New("Could not retrieve creation_date")
	}

	return creationDate, nil
}
