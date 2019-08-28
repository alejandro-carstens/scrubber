package runners

import (
	"errors"
	"regexp"
	"scrubber/actions/criterias"
	"scrubber/actions/infos"
	"sort"
	"time"

	"github.com/araddon/dateparse"
)

// NOT_VALID_DATE_PARAMETER error message
const NOT_VALID_DATE_PARAMETER = "Field is not a valid date parameter"

type useAgeFilterRunner struct {
	aggregateBaseRunner
}

func (uafr *useAgeFilterRunner) runAgeSorters(criteria criterias.Sortable) ([]string, error) {
	switch criteria.GetSource() {
	case "creation_date":
		return uafr.sortByCreationDate(criteria)
	case "name":
		return uafr.sortByName(criteria)
	case "field_stats":
		return uafr.sortByFieldStats(criteria)
	}

	return nil, errors.New("Invalid sortable criteria source")
}

func (uafr *useAgeFilterRunner) runDefaultSorter(criteria criterias.Sortable) []string {
	sortedList := []string{}

	for name, _ := range uafr.info {
		sortedList = append(sortedList, name)
	}

	if criteria.GetReverse() {
		sort.Sort(sort.Reverse(sort.StringSlice(sortedList)))
	} else {
		sort.Sort(sort.StringSlice(sortedList))
	}

	uafr.report.AddReason("Sorting by alphabetical order")
	uafr.report.AddReason("Sorted List: %v", sortedList)

	return sortedList
}

func (uafr *useAgeFilterRunner) sortByCreationDate(criteria criterias.Sortable) ([]string, error) {
	dateMap := map[time.Time]string{}
	timestamps := []time.Time{}

	for name, info := range uafr.info {
		creationDate, err := uafr.creantionDate(info)

		if err != nil {
			return nil, err
		}

		timestamp, err := time.Parse(time.RFC3339, creationDate)

		if err != nil {
			return nil, err
		}

		for {
			if _, valid := dateMap[timestamp]; !valid {
				dateMap[timestamp] = name
				break
			}

			timestamp = timestamp.Add(1 * time.Millisecond)
		}

		timestamps = append(timestamps, timestamp)
	}

	uafr.report.AddReason("Sorting by creation date")

	return uafr.createList(criteria, dateMap, timestamps)
}

func (uafr *useAgeFilterRunner) sortByName(criteria criterias.Sortable) ([]string, error) {
	var err error
	dateMap := map[time.Time]string{}
	timestamps := []time.Time{}

	switch criteria.GetTimestring() {
	case "Y.m.d":
		timestamps, dateMap, err = uafr.parseDatesFromName(`\d{4}.\d{2}.\d{2}`, criteria.GetStrictMode())
		break
	case "m.d.Y":
		timestamps, dateMap, err = uafr.parseDatesFromName(`\d{2}.\d{2}.\d{4}`, criteria.GetStrictMode())
		break
	case "Y.m":
		timestamps, dateMap, err = uafr.parseDatesFromName(`\d{4}.\d{2}`, criteria.GetStrictMode())
		break
	case "Y-m-d":
		timestamps, dateMap, err = uafr.parseDatesFromName(`\d{4}-\d{2}-\d{2}`, criteria.GetStrictMode())
		break
	case "Y-m-d H:M":
		timestamps, dateMap, err = uafr.parseDatesFromName(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}`, criteria.GetStrictMode())
		break
	case "Y-m-d H:M:S":
		timestamps, dateMap, err = uafr.parseDatesFromName(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`, criteria.GetStrictMode())
		break
	case "m-d-Y":
		timestamps, dateMap, err = uafr.parseDatesFromName(`\d{2}-\d{2}-\d{4}`, criteria.GetStrictMode())
		break
	case "Y-m":
		timestamps, dateMap, err = uafr.parseDatesFromName(`\d{4}-\d{2}`, criteria.GetStrictMode())
		break
	}

	if err != nil {
		return nil, err
	}

	uafr.report.AddReason("Sorting by name")

	return uafr.createList(criteria, dateMap, timestamps)
}

func (uafr *useAgeFilterRunner) sortByFieldStats(criteria criterias.Sortable) ([]string, error) {
	channel := make(chan *fieldStatsResponse, len(uafr.info))

	for name, _ := range uafr.info {
		go uafr.executeFieldStats(channel, name, criteria)
	}

	dateMap := map[time.Time]string{}
	timestamps := []time.Time{}

	for i := 0; i < len(uafr.info); i++ {
		fieldStatsResponse := <-channel

		if fieldStatsResponse.err != nil && fieldStatsResponse.err.Error() == NOT_VALID_DATE_PARAMETER && !criteria.GetStrictMode() {
			continue
		}

		if fieldStatsResponse.err != nil {
			return nil, fieldStatsResponse.err
		}

		timestamp := fieldStatsResponse.timestamp

		for {
			if _, valid := dateMap[timestamp]; !valid {
				dateMap[timestamp] = fieldStatsResponse.index

				break
			}

			timestamp = timestamp.Add(1 * time.Millisecond)
		}

		timestamps = append(timestamps, timestamp)
	}

	uafr.report.AddReason("Sorting by fields stats")

	return uafr.createList(criteria, dateMap, timestamps)
}

func (uafr *useAgeFilterRunner) parseDatesFromName(regPattern string, strictMode bool) ([]time.Time, map[time.Time]string, error) {
	dateMap := map[time.Time]string{}
	timestamps := []time.Time{}

	for name, _ := range uafr.info {
		reg, err := regexp.Compile(regPattern)

		if err != nil {
			return nil, nil, err
		}

		value := reg.FindString(name)

		if strictMode && len(value) == 0 {
			return nil, nil, errors.New("Could not match the expected pattern")
		} else if !strictMode && len(value) == 0 {
			uafr.report.AddReason(
				"Could not match the expected pattern '%v' for index/snapshot '%v'",
				regPattern,
				name,
			)

			continue
		}

		timestamp, err := dateparse.ParseLocal(value)

		if err != nil {
			return nil, nil, err
		}

		for {
			if _, valid := dateMap[timestamp]; !valid {
				dateMap[timestamp] = name

				break
			}

			timestamp.Add(1 * time.Millisecond)
		}

		timestamps = append(timestamps, timestamp)
	}

	return timestamps, dateMap, nil
}

func (uafr *useAgeFilterRunner) executeFieldStats(channel chan *fieldStatsResponse, index string, criteria criterias.Sortable) {
	response := new(fieldStatsResponse)

	result, err := uafr.connection.Builder(index).MinMax(criteria.GetField(), true)

	if err != nil {
		response.err = err
		channel <- response
		return
	}

	var date string
	var valid bool

	switch criteria.GetStatsResult() {
	case "min":
		date, valid = result.Min.(string)
		break
	case "max":
		date, valid = result.Max.(string)
		break
	}

	if !valid {
		response.err = errors.New(NOT_VALID_DATE_PARAMETER)
		channel <- response
		return
	}

	timestamp, err := dateparse.ParseLocal(date)

	if err != nil {
		response.err = err
		channel <- response
		return
	}

	response.index = index
	response.timestamp = timestamp
	channel <- response
}

func (uafr *useAgeFilterRunner) createList(criteria criterias.Sortable, indicesMap map[time.Time]string, timestamps []time.Time) ([]string, error) {
	var ts timeSlice = timestamps

	sort.Sort(ts)

	if criteria.GetReverse() {
		sort.Sort(sort.Reverse(ts))
	}

	indicesList := []string{}

	for _, key := range ts {
		index := indicesMap[key]

		uafr.report.AddReason("Index: %v, Timestamp: %v", index, key.Format(time.RFC3339))

		indicesList = append(indicesList, index)
	}

	return indicesList, nil
}

func (uafr *useAgeFilterRunner) creantionDate(info infos.Informable) (string, error) {
	creationDate := info.CreationDate()

	if len(creationDate) == 0 {
		return "", errors.New("Could not retrieve creation_date")
	}

	return creationDate, nil
}
