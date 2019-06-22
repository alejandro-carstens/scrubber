package tests

import (
	"os"
	"scrubber/console"
	"testing"

	"github.com/alejandro-carstens/golastic"
	"github.com/stretchr/testify/assert"
)

func TestSchedulerRunActionImmidiately(t *testing.T) {
	currentPath, err := os.Getwd()

	if err != nil {
		t.Error(err)
	}

	scheduler := console.NewScheduler(currentPath + "/../stubs/scheduler_extract_files_test/create_actions")

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
