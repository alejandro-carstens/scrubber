package tests

import (
	"sync"
	"testing"

	"github.com/alejandro-carstens/golastic"
	"github.com/stretchr/testify/assert"
)

func TestListSnapshots(t *testing.T) {
	var waitGroup sync.WaitGroup

	createIndexFiles := []string{
		"/test_files/create_index_by_name_1.yml",
		"/test_files/create_index_by_name_2.yml",
		"/test_files/create_index_by_name.yml",
	}

	waitGroup.Add(len(createIndexFiles))

	for _, createIndexFile := range createIndexFiles {
		go createTestIndexAsync(createIndexFile, &waitGroup)
	}

	waitGroup.Wait()

	takeAction("/test_files/create_repository.yml", t)

	createSnapshotFiles := []string{
		"/test_files/count_test_snapshot_index_1.yml",
		"/test_files/count_test_snapshot_index.yml",
	}

	for _, createSnapshotFile := range createSnapshotFiles {
		takeAction(createSnapshotFile, t)
	}

	action := takeAction("/test_files/list_snapshots.yml", t)

	expectedSnapshots := []string{
		"count_snapshot-2019.01.02",
		"count_snapshot-2019.01.01",
	}

	assert.ElementsMatch(t, expectedSnapshots, action.List())

	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		t.Error(err)
	}

	if _, err := builder.DeleteSnapshot("my_backup_repository", "count_snapshot-2019.01.01"); err != nil {
		t.Error(err)
	}

	if err := snapshotCleanup("my_backup_repository", "count_snapshot-2019.01.02", "_all", builder); err != nil {
		t.Error(err)
	}
}
