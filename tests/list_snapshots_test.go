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
		"/testfiles/create_index_by_name_1.yml",
		"/testfiles/create_index_by_name_2.yml",
		"/testfiles/create_index_by_name.yml",
	}

	waitGroup.Add(len(createIndexFiles))

	for _, createIndexFile := range createIndexFiles {
		go createTestIndexAsync(createIndexFile, &waitGroup)
	}

	waitGroup.Wait()

	takeAction("/testfiles/create_repository.yml", t)

	createSnapshotFiles := []string{
		"/testfiles/count_test_snapshot_index_1.yml",
		"/testfiles/count_test_snapshot_index.yml",
	}

	for _, createSnapshotFile := range createSnapshotFiles {
		takeAction(createSnapshotFile, t)
	}

	action := takeAction("/testfiles/list_snapshots.yml", t)

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
