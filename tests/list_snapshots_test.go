package tests

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListSnapshots(t *testing.T) {
	var waitGroup sync.WaitGroup

	createIndexFiles := []string{
		"/testdata/create_index_by_name_1.yml",
		"/testdata/create_index_by_name_2.yml",
		"/testdata/create_index_by_name.yml",
	}

	waitGroup.Add(len(createIndexFiles))

	for _, createIndexFile := range createIndexFiles {
		go createTestIndexAsync(createIndexFile, &waitGroup)
	}

	waitGroup.Wait()

	takeAction("/testdata/create_repository.yml", t)

	createSnapshotFiles := []string{
		"/testdata/count_test_snapshot_index_1.yml",
		"/testdata/count_test_snapshot_index.yml",
	}

	for _, createSnapshotFile := range createSnapshotFiles {
		takeAction(createSnapshotFile, t)
	}

	action := takeAction("/testdata/list_snapshots.yml", t)

	expectedSnapshots := []string{
		"count_snapshot-2019.01.02",
		"count_snapshot-2019.01.01",
	}

	assert.ElementsMatch(t, expectedSnapshots, action.List())

	connection := connection()

	if _, err := connection.Indexer(nil).DeleteSnapshot("my_backup_repository", "count_snapshot-2019.01.01"); err != nil {
		t.Error(err)
	}

	if err := snapshotCleanup("my_backup_repository", "count_snapshot-2019.01.02", "_all", connection); err != nil {
		t.Error(err)
	}
}
