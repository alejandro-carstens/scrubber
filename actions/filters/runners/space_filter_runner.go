package runners

import (
	"scrubber/actions/criterias"
	"scrubber/actions/responses"

	"github.com/alejandro-carstens/golastic"
)

const GIGABYTE_CONVERSION float64 = 1000000000
const MEGABYTE_CONVERSION float64 = 1000000

type IndexStatsResponse struct {
	sizeInBytes map[string]float64
	err         error
}

type spaceFilterRunner struct {
	useAgeFilterRunner
}

func (sfr *spaceFilterRunner) Init(info ...responses.Informable) (Runnerable, error) {
	if err := sfr.BaseInit(info...); err != nil {
		return nil, err
	}

	return sfr, nil
}

func (sfr *spaceFilterRunner) RunFilter(channel chan *FilterResponse, criteria criterias.Criteriable) {
	if err := sfr.validateCriteria(criteria); err != nil {
		channel <- sfr.response.setError(err)
		return
	}

	indicesStatsResponse := make(chan *IndexStatsResponse, 1)

	go sfr.executeIndexStats(indicesStatsResponse)

	var err error
	var sortedList []string

	space := criteria.(*criterias.Space)

	if space.UseAge {
		sortedList, err = sfr.runAgeSorters(space)
	} else {
		sortedList = sfr.runDefaultSorter(space)
	}

	statsResponse := <-indicesStatsResponse

	if statsResponse.err != nil {
		channel <- sfr.response.setError(statsResponse.err)
		return
	}

	sortedList = sfr.sumSpace(sortedList, statsResponse.sizeInBytes, space)

	if !space.Include() {
		sfr.report.AddReason("Excluding indices: '%v' from list", sortedList)

		sortedList = sfr.excludeIndices(sortedList)
	}

	sfr.report.AddResults(sortedList...)

	channel <- sfr.response.
		setError(err).
		setPassed(true).
		setReport(sfr.report).
		setList(sortedList)
}

func (sfr *spaceFilterRunner) executeIndexStats(indicesStatsResponse chan *IndexStatsResponse) {
	response := new(IndexStatsResponse)
	indices := []string{}

	for index, _ := range sfr.info {
		indices = append(indices, index)
	}

	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		response.err = err
		indicesStatsResponse <- response
		return
	}

	mapContainer, err := builder.IndexStats(indices...)

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

func (sfr *spaceFilterRunner) sumSpace(sortedList []string, sizeInBytes map[string]float64, space *criterias.Space) []string {
	count := float64(0)

	sfr.report.AddReason("Starting space sum")

	for key, index := range sortedList {
		count = count + sizeInBytes[index]

		convertedCount := sfr.convertCount(count, space.Units)

		sfr.report.AddReason(
			"Index %v space '%v' '%v', space sum '%v'",
			index,
			sfr.convertCount(sizeInBytes[index], space.Units),
			space.Units,
			convertedCount,
		)

		if space.FeThresholdBehavior == "greater_than" && convertedCount > float64(space.DiskSpace) {
			sfr.report.AddReason(
				"Threshold '%v' exceeded, space sum '%v'",
				space.DiskSpace,
				convertedCount,
			)

			return sortedList[:key+1]
		}
	}

	if space.FeThresholdBehavior == "less_than" && sfr.convertCount(count, space.Units) < float64(space.DiskSpace) {
		sfr.report.AddReason(
			"Threshold '%v' not exceeded, space sum '%v'",
			space.DiskSpace,
			sfr.convertCount(count, space.Units),
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
