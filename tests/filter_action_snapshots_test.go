package tests

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFilterSnapshots(t *testing.T) {
	for _, data := range filterSnapshotsDataProvider() {
		log.Println(fmt.Sprintf("Running %v", data["name"]))

		if _, err := createTestIndex(data["create_index"]); err != nil {
			t.Error(err)
		}

		time.Sleep(time.Duration(int64(2)) * time.Second)

		takeAction(data["create_repository"], t)
		takeAction(data["create_snapshot"], t)
		takeAction(data["action"], t)

		time.Sleep(time.Duration(int64(2)) * time.Second)

		connection := connection()

		list, err := connection.Indexer(nil).ListSnapshots(data["repository"])

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, 0, len(list))

		if err := snapshotCleanup(data["repository"], "", data["index_name"], connection); err != nil {
			t.Error(err)
		}

		log.Println(fmt.Sprintf("Done running %v", data["name"]))
	}
}

func TestFilterSnapshotsByCount(t *testing.T) {
	for _, data := range filterSnapshotsByCountDataProvider() {
		log.Println("Running " + data["name"])

		var waitGroup sync.WaitGroup

		createIndexFiles := strings.Split(data["create_indices"], ",")

		if len(createIndexFiles) == 0 {
			t.Error(errors.New("There most be at least one index for creation"))
		}

		waitGroup.Add(len(createIndexFiles))

		for _, createIndexFile := range createIndexFiles {
			go createTestIndexAsync(createIndexFile, &waitGroup)
		}

		waitGroup.Wait()

		takeAction(data["create_repository"], t)

		createSnapshotFiles := strings.Split(data["create_snapshots"], ",")

		if len(createSnapshotFiles) == 0 {
			t.Error(errors.New("There most be at least one snapshot to create"))
		}

		waitGroup.Add(len(createSnapshotFiles))

		for i, createSnapshotFile := range createSnapshotFiles {
			go takeActionAsync(createSnapshotFile, t, &waitGroup)

			if _, valid := data["wait"]; valid && i+1 < len(createSnapshotFiles) {
				time.Sleep(time.Duration(int64(2)) * time.Second)
			}
		}

		waitGroup.Wait()

		takeAction(data["action"], t)

		time.Sleep(time.Duration(int64(2)) * time.Second)

		connection := connection()

		list, err := connection.Indexer(nil).ListSnapshots(data["repository"])

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, data["expected_snapshot_count"], fmt.Sprint(len(list)))
		assert.Equal(t, data["existing_snapshot"], list[0])

		if err := snapshotCleanup(data["repository"], data["existing_snapshot"], data["index_name"], connection); err != nil {
			t.Error(err)
		}

		log.Println("Done running " + data["name"])
	}
}

func filterSnapshotsDataProvider() []map[string]string {
	data := []map[string]string{}

	data = append(data, map[string]string{
		"name":              "TestFilterByCreationDate",
		"create_index":      "/testdata/create_pattern_index.yml",
		"create_repository": "/testdata/create_repository.yml",
		"create_snapshot":   "/testdata/snapshot_index.yml",
		"action":            "/testdata/delete_snapshot_by_creation_date.yml",
		"repository":        "my_backup_repository",
		"index_name":        "alejandro-carstens-1992.06.02",
	})
	data = append(data, map[string]string{
		"name":              "TestFilterByName",
		"create_index":      "/testdata/create_pattern_index.yml",
		"create_repository": "/testdata/create_repository.yml",
		"create_snapshot":   "/testdata/pattern_snapshot.yml",
		"action":            "/testdata/delete_snapshot_by_name.yml",
		"repository":        "my_backup_repository",
		"index_name":        "alejandro-carstens-1992.06.02",
	})
	data = append(data, map[string]string{
		"name":              "TestFilterByPattern",
		"create_index":      "/testdata/create_pattern_index.yml",
		"create_repository": "/testdata/create_repository.yml",
		"create_snapshot":   "/testdata/pattern_snapshot.yml",
		"action":            "/testdata/delete_snapshot_by_pattern.yml",
		"repository":        "my_backup_repository",
		"index_name":        "alejandro-carstens-1992.06.02",
	})

	return data
}

func filterSnapshotsByCountDataProvider() []map[string]string {
	data := []map[string]string{}

	data = append(data, map[string]string{
		"name":                    "TestFilterByCountSortByCreationDate",
		"create_indices":          "/testdata/create_index_by_name.yml,/testdata/create_index_by_name_1.yml,/testdata/create_index_by_name_2.yml",
		"create_repository":       "/testdata/create_repository.yml",
		"create_snapshots":        "/testdata/count_test_snapshot_index_1.yml,/testdata/count_test_snapshot_index.yml",
		"action":                  "/testdata/delete_snapshot_by_count_sort_by_creation_date.yml",
		"repository":              "my_backup_repository",
		"index_name":              "_all",
		"wait":                    "true",
		"existing_snapshot":       "count_snapshot-2019.01.01",
		"expected_snapshot_count": "1",
	})
	data = append(data, map[string]string{
		"name":                    "TestFilterByCountSortByName",
		"create_indices":          "/testdata/create_index_by_name.yml,/testdata/create_index_by_name_1.yml,/testdata/create_index_by_name_2.yml",
		"create_repository":       "/testdata/create_repository.yml",
		"create_snapshots":        "/testdata/count_test_snapshot_index_1.yml,/testdata/count_test_snapshot_index.yml",
		"action":                  "/testdata/delete_snapshot_by_count_sort_by_name.yml",
		"repository":              "my_backup_repository",
		"index_name":              "_all",
		"wait":                    "true",
		"existing_snapshot":       "count_snapshot-2019.01.02",
		"expected_snapshot_count": "1",
	})
	data = append(data, map[string]string{
		"name":                    "TestFilterByCountSortByPattern",
		"create_indices":          "/testdata/create_index_by_name.yml,/testdata/create_index_by_name_1.yml,/testdata/create_index_by_name_2.yml",
		"create_repository":       "/testdata/create_repository.yml",
		"create_snapshots":        "/testdata/count_test_snapshot_index_1.yml,/testdata/count_test_snapshot_index.yml",
		"action":                  "/testdata/delete_snapshot_by_count_sort_by_pattern.yml",
		"repository":              "my_backup_repository",
		"index_name":              "_all",
		"wait":                    "true",
		"existing_snapshot":       "count_snapshot-2019.01.01",
		"expected_snapshot_count": "1",
	})

	return data
}
