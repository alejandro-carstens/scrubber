package runners

import (
	"errors"
	"scrubber/actions/infos"
	"time"

	"github.com/alejandro-carstens/golastic"
)

const SECONDS_PER_DAY float64 = 86400
const SECONDS_PER_MONTH float64 = 2628000
const SECONDS_PER_YEAR float64 = 31540000

type fieldStatsResponse struct {
	err       error
	timestamp time.Time
	index     string
}

type timeSlice []time.Time

// Less timeSlice sorting method
func (ts timeSlice) Less(i int, j int) bool {
	return ts[i].Before(ts[j])
}

// Swap timeSlice sorting method
func (ts timeSlice) Swap(i int, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

// Len timeSlice sorting method
func (ts timeSlice) Len() int {
	return len(ts)
}

// NewRunner return a filter runner
func NewRunner(criteria string, builder *golastic.ElasticsearchBuilder, info ...infos.Informable) (Runnerable, error) {
	var runner Runnerable

	switch criteria {
	case "age":
		runner = new(ageFilterRunner)
		break
	case "pattern":
		runner = new(patternFilterRunner)
		break
	case "count":
		runner = new(countFilterRunner)
		break
	case "closed":
		runner = new(closedFilterRunner)
		break
	case "empty":
		runner = new(emptyFilterRunner)
		break
	case "kibana":
		runner = new(kibanaFilterRunner)
		break
	case "alias":
		runner = new(aliasFilterRunner)
		break
	case "allocated":
		runner = new(allocatedFilterRunner)
		break
	case "forcemerged":
		runner = new(forcemergedFilterRunner)
		break
	case "space":
		runner = new(spaceFilterRunner)
		break
	case "state":
		runner = new(stateFilterRunner)
	default:
		return nil, errors.New("Invalid criteria")
	}

	return runner.Init(builder, info...)
}

func elapsed(from, to time.Time) (years, months, days, hours, minutes, seconds int) {
	if from.Location() != to.Location() {
		to = to.In(from.Location())
	}

	if from.After(to) {
		from, to = to, from
	}

	y1, M1, d1 := from.Date()
	y2, M2, d2 := to.Date()

	h1, m1, s1 := from.Clock()
	h2, m2, s2 := to.Clock()

	years = y2 - y1
	months = int(M2 - M1)
	days = d2 - d1

	hours = h2 - h1
	minutes = m2 - m1
	seconds = s2 - s1

	if seconds < 0 {
		seconds += 60
		minutes--
	}

	if minutes < 0 {
		minutes += 60
		hours--
	}

	if hours < 0 {
		hours += 24
		days--
	}

	if days < 0 {
		days += daysIn(y2, M2-1)
		months--
	}

	if months < 0 {
		months += 12
		years--
	}

	return
}

func daysIn(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func secondsToDays(seconds float64) float64 {
	return seconds / SECONDS_PER_DAY
}

func secondsToMonths(seconds float64) float64 {
	return seconds / SECONDS_PER_MONTH
}

func secondsToYears(seconds float64) float64 {
	return seconds / SECONDS_PER_YEAR
}
