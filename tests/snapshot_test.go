package tests

import (
	"testing"
	"time"
)

func TestSnapshot(t *testing.T) {
	if _, err := createTestIndex("/testdata/create_pattern_index.yml"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	takeAction("/testdata/create_repository.yml", t)

	time.Sleep(time.Duration(int64(3)) * time.Second)

	takeAction("/testdata/snapshot_index.yml", t)

	if err := snapshotCleanup("my_backup_repository", "my_first_snapshot", "alejandro-carstens-1992.06.02", nil); err != nil {
		t.Error(err)
	}
}
