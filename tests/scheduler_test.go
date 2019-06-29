package tests

import (
	"os"
	"scrubber/console"
	"scrubber/logging"
	"testing"

	"github.com/alejandro-carstens/golastic"
	"github.com/stretchr/testify/assert"
)

func TestSchedulerRunActionImmidiately(t *testing.T) {
	currentPath, err := os.Getwd()

	if err != nil {
		t.Error(err)
	}

	filePath := currentPath + "/test_files/scheduler_extract_files_test/create_actions"
	logger := logging.NewSrvLogger("", true, true, true, true)

	scheduler := console.NewScheduler(filePath, logger)

	if err := scheduler.Run(); err != nil {
		t.Error(err)
	}

	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		t.Error(err)
	}

	exists, err := builder.Exists("my_index")

	if err != nil {
		t.Error(err)
	}

	assert.True(t, exists)

	if err := builder.DeleteIndex("my_index"); err != nil {
		t.Error(err)
	}

	exists, err = builder.Exists("my_async_index")

	if err != nil {
		t.Error(err)
	}

	assert.True(t, exists)

	if err := builder.DeleteIndex("my_async_index"); err != nil {
		t.Error(err)
	}
}
