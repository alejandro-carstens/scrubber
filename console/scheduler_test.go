package console

import (
	"os"
	rp "scrubber/resourcepool"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchedulerExtractFiles(t *testing.T) {
	rp.Boot("mysql")

	currentPath, err := os.Getwd()

	if err != nil {
		t.Error(err)
	}

	filePath := currentPath + "/../tests/testdata/schedulerdata"

	configs, err := NewScheduler(filePath, []string{}, nil).extractConfigs()

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
	rp.Boot("mysql")

	currentPath, err := os.Getwd()

	if err != nil {
		t.Error(err)
	}

	filePath := currentPath + "/../tests/testdata/schedulerdata/deleteactions/aggregate"

	scheduler := NewScheduler(filePath, []string{}, nil)

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
