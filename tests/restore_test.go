package tests

import (
	"testing"
	"time"

	"github.com/alejandro-carstens/golastic"
	"github.com/stretchr/testify/assert"
)

func TestRestoreSnapshot(t *testing.T) {
	for _, data := range restoreSnapshotDataProvider() {
		if _, err := createTestIndex(data["create_index"]); err != nil {
			t.Error(err)
		}

		time.Sleep(time.Duration(int64(2)) * time.Second)

		takeAction(data["create_repository"], t)
		takeAction(data["create_snapshot"], t)
		takeAction(data["restore_snapshot"], t)

		time.Sleep(time.Duration(int64(2)) * time.Second)

		builder, err := golastic.NewBuilder(nil, nil)

		if err != nil {
			t.Error(err)
		}

		list, err := builder.ListSnapshots(data["repository"])

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, 1, len(list))

		if err := builder.DeleteIndex(data["index_name"]); err != nil {
			t.Error(err)
		}

		if err := snapshotCleanup(data["repository"], data["snapshot"], data["restored_index"], builder); err != nil {
			t.Error(err)
		}
	}
}

func restoreSnapshotDataProvider() []map[string]string {
	data := []map[string]string{}

	data = append(data, map[string]string{
		"create_index":      "/../stubs/create_index_to_restore.yml",
		"create_repository": "/../stubs/create_repository.yml",
		"create_snapshot":   "/../stubs/snapshot_to_restore.yml",
		"restore_snapshot":  "/../stubs/restore_snapshot_by_creation_date.yml",
		"repository":        "my_backup_repository",
		"snapshot":          "my_first_snapshot",
		"restored_index":    "restored_index_1",
		"index_name":        "index_1",
	})
	data = append(data, map[string]string{
		"create_index":      "/../stubs/create_index_to_restore.yml",
		"create_repository": "/../stubs/create_repository.yml",
		"create_snapshot":   "/../stubs/snapshot_to_restore.yml",
		"restore_snapshot":  "/../stubs/restore_snapshot_by_creation_date_no_wait.yml",
		"repository":        "my_backup_repository",
		"snapshot":          "my_first_snapshot",
		"restored_index":    "restored_index_1",
		"index_name":        "index_1",
	})

	return data
}
