package runners

import (
	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/criterias"
	"github.com/alejandro-carstens/scrubber/actions/infos"
)

const GIGABYTE_CONVERSION float64 = 1000000000
const MEGABYTE_CONVERSION float64 = 1000000

// IndexStatsResponse represents the response of an index stats call
type IndexStatsResponse struct {
	sizeInBytes map[string]float64
	err         error
}

type spaceFilterRunner struct {
	useAgeFilterRunner
	criteria *criterias.Space
}

// Init initializes the filter runner
func (sfr *spaceFilterRunner) Init(criteria criterias.Criteriable, builder *golastic.Connection, info ...infos.Informable) (Runnerable, error) {
	if err := sfr.BaseInit(criteria, builder, info...); err != nil {
		return nil, err
	}

	sfr.criteria = criteria.(*criterias.Space)

	return sfr, nil
}

// RunFilter filters out elements from the actionable list
func (sfr *spaceFilterRunner) RunFilter(channel chan *FilterResponse) {
	indicesStatsResponse := make(chan *IndexStatsResponse, 1)

	go sfr.executeIndexStats(indicesStatsResponse)

	var err error
	var sortedList []string

	if sfr.criteria.UseAge {
		sortedList, err = sfr.runAgeSorters(sfr.criteria)
	} else {
		sortedList = sfr.runDefaultSorter(sfr.criteria)
	}

	statsResponse := <-indicesStatsResponse

	if statsResponse.err != nil {
		channel <- sfr.response.setError(statsResponse.err)
		return
	}

	sortedList = sfr.sumSpace(sortedList, statsResponse.sizeInBytes)

	if !sfr.criteria.Include() {
		sfr.report.AddReason("Excluding indices: '%v' from list", sortedList)

		sortedList = sfr.excludeIndices(sortedList)
	}

	sfr.report.AddResults(sortedList...)

	channel <- sfr.response.setError(err).setPassed(true).setReport(sfr.report).setList(sortedList)
}

func (sfr *spaceFilterRunner) executeIndexStats(indicesStatsResponse chan *IndexStatsResponse) {
	response := new(IndexStatsResponse)
	indices := []string{}

	for index, _ := range sfr.info {
		indices = append(indices, index)
	}

	mapContainer, err := sfr.connection.Indexer(nil).IndexStats(indices...)

	if err != nil {
		response.err = err
		indicesStatsResponse <- response
		return
	}

	sizeInBytesMap := map[string]float64{}

	for index, child := range mapContainer {
		sizeInBytes, valid := child.S("total", "store", "size_in_bytes").Data().(float64)

		if !valid {
			response.err = err
			break
		}

		sizeInBytesMap[index] = sizeInBytes
	}

	response.sizeInBytes = sizeInBytesMap

	indicesStatsResponse <- response
}

func (sfr *spaceFilterRunner) sumSpace(sortedList []string, sizeInBytes map[string]float64) []string {
	count := float64(0)

	sfr.report.AddReason("Starting space sum")

	for key, index := range sortedList {
		count = count + sizeInBytes[index]

		convertedCount := sfr.convertCount(count, sfr.criteria.Units)

		sfr.report.AddReason(
			"Index %v space '%v' '%v', space sum '%v'",
			index,
			sfr.convertCount(sizeInBytes[index], sfr.criteria.Units),
			sfr.criteria.Units,
			convertedCount,
		)

		if sfr.criteria.FeThresholdBehavior == "greater_than" && convertedCount > float64(sfr.criteria.DiskSpace) {
			sfr.report.AddReason(
				"Threshold '%v' exceeded, space sum '%v'",
				sfr.criteria.DiskSpace,
				convertedCount,
			)

			return sortedList[:key+1]
		}
	}

	if sfr.criteria.FeThresholdBehavior == "less_than" && sfr.convertCount(count, sfr.criteria.Units) < float64(sfr.criteria.DiskSpace) {
		sfr.report.AddReason(
			"Threshold '%v' not exceeded, space sum '%v'",
			sfr.criteria.DiskSpace,
			sfr.convertCount(count, sfr.criteria.Units),
		)

		return sortedList
	}

	return []string{}
}

func (sfr *spaceFilterRunner) convertCount(count float64, units string) float64 {
	if units == "MB" {
		return count / MEGABYTE_CONVERSION
	}

	return count / GIGABYTE_CONVERSION
}
