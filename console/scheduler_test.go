package console

import (
	"os"
	"scrubber/logger"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchedulerExtractFiles(t *testing.T) {
	currentPath, err := os.Getwd()

	if err != nil {
		t.Error(err)
	}

	filePath := currentPath + "/../tests/test_files/scheduler_extract_files_test"
	logger := logger.NewLogger("", true, true, true, true)

	configs, err := NewScheduler(filePath, logger).extractConfigs()

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 4, len(configs))

	for path, config := range configs {
		description := config.S("description").Data().(string)

		if strings.Contains(path, "delete_indices_by_alias.yml") {
			assert.Equal(t, "Deletes indices by alias", description)
		} else if strings.Contains(path, "create_index.yml") {
			assert.Equal(t, "Creates the specified index with the specified settings", description)
		} else if strings.Contains(path, "create_index_async.yml") {
			assert.Equal(t, "Creates the specified index with the specified settings", description)
		} else {
			assert.Equal(t, "Delete 3 indices based on space order alphabetically", description)
		}
	}
}

func TestSchedulerSchedule(t *testing.T) {
	currentPath, err := os.Getwd()

	if err != nil {
		t.Error(err)
	}

	filePath := currentPath + "/../tests/test_files/scheduler_extract_files_test/delete_actions/aggregate"
	logger := logger.NewLogger("", true, true, true, true)

	scheduler := NewScheduler(filePath, logger)

	configs, err := scheduler.extractConfigs()

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1, len(configs))

	for _, config := range configs {
		task, err := scheduler.schedule(config)

		if err != nil {
			t.Error(err)
		}

		assert.NotEmpty(t, task)
	}
}
