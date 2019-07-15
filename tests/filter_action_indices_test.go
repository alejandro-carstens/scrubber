package tests

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/golastic"
	"github.com/stretchr/testify/assert"
)

const SEED_INDEX_COUNT int = 4
const ELASTICSEARCH_BULK_INSERT_LIMIT int = 10000

func TestFilterIndices(t *testing.T) {
	for _, data := range filterIndicesDataProvider() {
		log.Println("Running " + data["test_name"])

		if _, err := createTestIndex(data["create_mock"]); err != nil {
			t.Error(err)
		}

		sleepTime := int64(2)

		if _, valid := data["sleep_time"]; valid {
			extraTime, err := strconv.Atoi(data["sleep_time"])

			if err != nil {
				t.Error(err)
			}

			sleepTime = sleepTime + int64(extraTime)
		}

		time.Sleep(time.Duration(sleepTime) * time.Second)

		takeAction(data["action_mock"], t)

		builder, err := golastic.NewBuilder(nil, nil)

		if err != nil {
			t.Error(err)
		}

		exists, err := builder.Exists(data["index_name"])

		if err != nil {
			t.Error(err)
		}

		assert.False(t, exists)

		log.Println("Done running " + data["test_name"])
	}
}

func TestFilterInidicesByFieldStats(t *testing.T) {
	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		t.Error(err)
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	go seedIndexAsync("my_index", 10, builder, &waitGroup, false)

	waitGroup.Wait()

	takeAction("/testdata/delete_indices_by_field_stats.yml", t)

	exists, err := builder.Exists("my_index")

	if err != nil {
		t.Error(err)
	}

	assert.False(t, exists)
}

func TestFilterClosedIndex(t *testing.T) {
	if _, err := createTestIndex("/testdata/create_index.yml"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		t.Error(err)
	}

	if _, err := builder.Close("my_index"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	takeAction("/testdata/delete_closed_index.yml", t)

	exists, err := builder.Exists("my_index")

	if err != nil {
		t.Error(err)
	}

	assert.False(t, exists)
}

func TestNoFilters(t *testing.T) {
	if _, err := createTestIndex("/testdata/create_index.yml"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	takeAction("/testdata/delete_indices_no_filters.yml", t)

	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		t.Error(err)
	}

	exists, err := builder.Exists("my_index")

	if err != nil {
		t.Error(err)
	}

	assert.False(t, exists)
}

func TestFilterIndicesByAlias(t *testing.T) {
	if _, err := createTestIndex("/testdata/create_index.yml"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		t.Error(err)
	}

	if _, err := builder.AddAlias("my_index", "my_alias"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	takeAction("/testdata/delete_indices_by_alias.yml", t)

	exists, err := builder.Exists("my_index")

	if err != nil {
		t.Error(err)
	}

	assert.False(t, exists)
}

func TestFilterIndicesByCount(t *testing.T) {
	for _, data := range filterIndicesByCountDataProvider() {
		log.Println("Running " + data["test_name"])

		var waitGroup sync.WaitGroup

		createIndexFiles := strings.Split(data["create_mocks"], ",")

		if len(createIndexFiles) == 0 {
			t.Error(errors.New("There most be at least one index for creation"))
		}

		waitGroup.Add(len(createIndexFiles))

		waitTime := 0

		res, valid := data["wait_time"]

		if valid {
			res, err := strconv.Atoi(res)

			if err != nil {
				t.Error(err)
			}

			waitTime = res
		}

		for _, createIndexFile := range createIndexFiles {
			if waitTime > 0 {
				time.Sleep(time.Duration(int64(waitTime)) * time.Second)
			}

			go createTestIndexAsync(createIndexFile, &waitGroup)
		}

		waitGroup.Wait()

		builder, err := golastic.NewBuilder(nil, nil)

		if err != nil {
			t.Error(err)
		}

		list, err := builder.ListIndices()

		if err != nil {
			t.Error(err)
		}

		if len(list) != len(createIndexFiles) {
			t.Error(errors.New("The number of indices created must match the number of indices passed in"))
		}

		takeAction(data["action_mock"], t)

		resultList, err := builder.ListIndices()

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, data["expected_index_count"], fmt.Sprint(len(resultList)))

		expectedExistingIndex, valid := data["expected_existing_index"]

		if valid {
			exists, err := builder.Exists(expectedExistingIndex)

			if err != nil {
				t.Error(err)
			}

			assert.True(t, exists)

			if err := builder.DeleteIndex(expectedExistingIndex); err != nil {
				t.Error(err)
			}
		}

		log.Println("Done running " + data["test_name"])
	}
}

func TestFilterIndicesByCountSortedByFieldStats(t *testing.T) {
	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		t.Error(err)
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(3)

	for i := 1; i <= 3; i++ {
		index := fmt.Sprint("my_index_" + fmt.Sprint(i))

		go seedIndexAsync(index, 100*i, builder, &waitGroup, false)
	}

	waitGroup.Wait()

	takeAction("/testdata/delete_indices_by_count_sorted_by_field_stats.yml", t)

	resultList, err := builder.ListIndices()

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "1", fmt.Sprint(len(resultList)))

	exists, err := builder.Exists("my_index_3")

	if err != nil {
		t.Error(err)
	}

	assert.True(t, exists)

	if err := builder.DeleteIndex("my_index_3"); err != nil {
		t.Error(err)
	}
}

func TestFilterIndicesBySpace(t *testing.T) {
	for _, data := range filterIndicesBySpaceDataProvider() {
		log.Println("Running " + data["test_name"])

		builder, err := golastic.NewBuilder(nil, nil)

		if err != nil {
			t.Error(err)
		}

		indices, valid := data["indices"]

		if !valid {
			t.Error(errors.New("No indices indices provided"))
		}

		indexSlice := strings.Split(indices, ",")

		var waitGroup sync.WaitGroup
		waitGroup.Add(len(indexSlice))

		for i, index := range indexSlice {
			if i >= 1 {
				i = SEED_INDEX_COUNT
			}

			go seedIndexAsync(index, 10000*(i+1), builder, &waitGroup, true)
		}

		waitGroup.Wait()

		actionMock, valid := data["action_mock"]

		if !valid {
			t.Error(errors.New("No action mock provided"))
		}

		config, err := config(actionMock)

		if err != nil {
			t.Error(err)
		}

		action, err := getAction(config)

		if err != nil {
			t.Error(err)
		}

		stats, err := builder.IndexStats(indexSlice...)

		if err != nil {
			t.Error(err)
		}

		assert.False(t, action.Perform().HasErrors())

		expectedIndexCount, valid := data["expected_index_count"]

		if valid {
			resultList, err := builder.ListIndices()

			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, expectedIndexCount, fmt.Sprint(len(resultList)))
		} else {
			assertSpaceFilteredIndices(t, data, stats, builder)
		}

		log.Println("Done running " + data["test_name"])
	}
}

func TestFilterIndicesByAllocation(t *testing.T) {
	if _, err := createTestIndex("/testdata/create_index.yml"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		t.Error(err)
	}

	body := `{"index.routing.allocation.require.shards": 1}`

	if _, err := builder.PutSettings(body, "my_index"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	takeAction("/testdata/delete_indices_by_allocation.yml", t)

	exists, err := builder.Exists("my_index")

	if err != nil {
		t.Error(err)
	}

	assert.False(t, exists)
}

func TestFilterIndicesByForcemerged(t *testing.T) {
	if _, err := createTestIndex("/testdata/create_index.yml"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		t.Error(err)
	}

	takeAction("/testdata/delete_indices_by_forcemerged.yml", t)

	exists, err := builder.Exists("my_index")

	if err != nil {
		t.Error(err)
	}

	assert.False(t, exists)
}

func filterIndicesDataProvider() []map[string]string {
	dataProvider := []map[string]string{}
	dataProvider = append(dataProvider, map[string]string{
		"test_name":   "TestDeleteIndicesByCreationDate",
		"create_mock": "/testdata/create_index.yml",
		"action_mock": "/testdata/delete_indices_by_creation_date.yml",
		"index_name":  "my_index",
		"sleep_time":  "1",
	})
	dataProvider = append(dataProvider, map[string]string{
		"test_name":   "TestDeleteIndicesByName",
		"create_mock": "/testdata/create_index_by_name.yml",
		"action_mock": "/testdata/delete_indices_by_name.yml",
		"index_name":  "my_index-2019-01-01",
	})
	dataProvider = append(dataProvider, map[string]string{
		"test_name":   "TestDeleteEmptyIndex",
		"create_mock": "/testdata/create_index.yml",
		"action_mock": "/testdata/delete_empty_index.yml",
		"index_name":  "my_index",
	})
	dataProvider = append(dataProvider, map[string]string{
		"test_name":   "TestDeleteKibanaIndex",
		"create_mock": "/testdata/create_kibana_index.yml",
		"action_mock": "/testdata/delete_kibana_index.yml",
		"index_name":  ".kibana",
	})
	dataProvider = append(dataProvider, map[string]string{
		"test_name":   "TestDeleteIndexByRegexPattern",
		"create_mock": "/testdata/create_pattern_index.yml",
		"action_mock": "/testdata/delete_regex_pattern_index.yml",
		"index_name":  "alejandro-carstens-1992.06.02",
	})
	dataProvider = append(dataProvider, map[string]string{
		"test_name":   "TestDeleteIndexByPrefixPattern",
		"create_mock": "/testdata/create_pattern_index.yml",
		"action_mock": "/testdata/delete_prefix_pattern_index.yml",
		"index_name":  "alejandro-carstens-1992.06.02",
	})
	dataProvider = append(dataProvider, map[string]string{
		"test_name":   "TestDeleteIndexBySuffixPattern",
		"create_mock": "/testdata/create_pattern_index.yml",
		"action_mock": "/testdata/delete_suffix_pattern_index.yml",
		"index_name":  "alejandro-carstens-1992.06.02",
	})
	dataProvider = append(dataProvider, map[string]string{
		"test_name":   "TestDeleteIndexByTimestringPattern",
		"create_mock": "/testdata/create_pattern_index.yml",
		"action_mock": "/testdata/delete_timestring_pattern_index.yml",
		"index_name":  "alejandro-carstens-1992.06.02",
	})

	return dataProvider
}

func filterIndicesByCountDataProvider() []map[string]string {
	dataProvider := []map[string]string{}

	dataProvider = append(dataProvider, map[string]string{
		"test_name":               "TestDeleteIndicesByCountSortedByTimestring",
		"create_mocks":            "/testdata/create_index_by_name.yml,/testdata/create_index_by_name_1.yml,/testdata/create_index_by_name_2.yml",
		"action_mock":             "/testdata/delete_indices_by_count_sorted_by_timestring.yml",
		"expected_index_count":    "1",
		"expected_existing_index": "my_index-2019-01-01",
	})
	dataProvider = append(dataProvider, map[string]string{
		"test_name":               "TestDeleteIndicesByCountNoSortingParam",
		"create_mocks":            "/testdata/create_index_by_name.yml,/testdata/create_index_by_name_1.yml,/testdata/create_index_by_name_2.yml",
		"action_mock":             "/testdata/delete_indices_by_count_no_sorting_param.yml",
		"expected_index_count":    "1",
		"expected_existing_index": "my_index-2019-01-03",
	})
	dataProvider = append(dataProvider, map[string]string{
		"test_name":               "TestDeleteIndicesByCountSortByPattern",
		"create_mocks":            "/testdata/create_index_by_name.yml,/testdata/create_pattern_index.yml,/testdata/create_pattern_index_1.yml",
		"action_mock":             "/testdata/delete_indices_by_count_sort_by_pattern.yml",
		"expected_index_count":    "1",
		"expected_existing_index": "my_index-2019-01-01",
	})
	dataProvider = append(dataProvider, map[string]string{
		"test_name":               "TestDeleteIndicesByCountSortByCreationDate",
		"create_mocks":            "/testdata/create_index_by_name.yml,/testdata/create_index_by_name_1.yml,/testdata/create_index_by_name_2.yml",
		"action_mock":             "/testdata/delete_indices_by_count_sort_by_creation_date.yml",
		"expected_index_count":    "1",
		"expected_existing_index": "my_index-2019-01-01",
		"wait_time":               "3",
	})

	return dataProvider
}

func filterIndicesBySpaceDataProvider() []map[string]string {
	dataProvider := []map[string]string{}

	dataProvider = append(dataProvider, map[string]string{
		"test_name":       "TestDeleteIndicesBySpaceSortedByTimestring",
		"indices":         "my_index-2017-03-01,my_index-2017-03-02,my_index-2017-03-03",
		"action_mock":     "/testdata/delete_indices_by_space_sorted_by_timestring.yml",
		"disk_space":      "1",
		"ordered_indices": "my_index-2017-03-03,my_index-2017-03-02,my_index-2017-03-01",
	})
	dataProvider = append(dataProvider, map[string]string{
		"test_name":            "TestDeleteIndicesBySpaceSortedLessThan",
		"indices":              "my_index-2017-03-07,my_index-2017-03-08,my_index-2017-03-09",
		"action_mock":          "/testdata/delete_indices_by_space_less_than.yml",
		"expected_index_count": "0",
	})
	dataProvider = append(dataProvider, map[string]string{
		"test_name":       "TestDeleteIndicesBySpaceSortedAlphabetically",
		"indices":         "my_index-2017-03-06,my_index-2017-03-04,my_index-2017-03-05",
		"action_mock":     "/testdata/delete_indices_by_space_sorted_alphabetically.yml",
		"disk_space":      "1",
		"ordered_indices": "my_index-2017-03-04,my_index-2017-03-05,my_index-2017-03-06",
	})

	return dataProvider
}

func seedIndexAsync(index string, count int, builder *golastic.ElasticsearchBuilder, waitGroup *sync.WaitGroup, useConstantTime bool) {
	defer waitGroup.Done()

	inserts := []golastic.ElasticModelable{}

	for i := 0; i < count; i++ {
		insert := golastic.NewGolasticModel()

		value := map[string]interface{}{}

		value["id"] = strconv.Itoa(i + 1 + 1000000000)
		value["exception"] = "Exception exception exception exception exception exception exception exception exception"
		value["request"] = "Request request request request request request request request request"
		value["message"] = "Message message message message message message message message message"
		value["bytes"] = int64(i)
		value["number"] = float64(i)

		if useConstantTime {
			constantTime, err := time.Parse(time.RFC3339, "2017-11-12T11:45:26.371Z")

			if err != nil {
				panic(err)
			}

			value["created_at"] = constantTime
		} else {
			value["created_at"] = time.Now().Add(time.Duration(int64(-1*(i+1))) * time.Hour)
		}

		insert.SetData(value)
		insert.SetIndex(index)

		inserts = append(inserts, insert)

		if count >= ELASTICSEARCH_BULK_INSERT_LIMIT && (i+1)%ELASTICSEARCH_BULK_INSERT_LIMIT == 0 {
			if _, err := builder.Insert(inserts...); err != nil {
				panic(err)
			}

			inserts = []golastic.ElasticModelable{}
		}
	}

	if count < ELASTICSEARCH_BULK_INSERT_LIMIT {
		if _, err := builder.Insert(inserts...); err != nil {
			panic(err)
		}
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)
}

func assertSpaceFilteredIndices(t *testing.T, data map[string]string, stats map[string]*gabs.Container, builder golastic.Queryable) {
	expectedExistingIndices := []string{}

	orderedIndices := strings.Split(data["ordered_indices"], ",")

	if len(orderedIndices) == 0 {
		t.Error(errors.New("ordered_indices must be specified"))
	}

	diskSpace, err := strconv.Atoi(data["disk_space"])

	if err != nil {
		t.Error(errors.New("disk_space must be specified with ordered_indices"))
	}

	count := float64(0)

	for i, index := range orderedIndices {
		count = count + stats[index].S("total", "store", "size_in_bytes").Data().(float64)

		if count/1000000 > float64(diskSpace) {
			expectedExistingIndices = orderedIndices[i+1:]
			break
		}
	}

	resultList, err := builder.ListIndices()

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(expectedExistingIndices), len(resultList))

	for _, expectedExistingIndex := range expectedExistingIndices {
		exists, err := builder.Exists(expectedExistingIndex)

		if err != nil {
			t.Error(err)
		}

		assert.True(t, exists)

		if err := builder.DeleteIndex(expectedExistingIndex); err != nil {
			t.Error(err)
		}
	}
}
