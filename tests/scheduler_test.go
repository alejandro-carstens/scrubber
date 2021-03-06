package tests

import (
	"context"
	"os"
	"testing"

	"github.com/alejandro-carstens/scrubber/console"
	"github.com/alejandro-carstens/scrubber/logger"
	"github.com/stretchr/testify/assert"
)

func TestSchedulerRunActionImmidiately(t *testing.T) {
	currentPath, err := os.Getwd()

	if err != nil {
		t.Error(err)
	}

	filePath := currentPath + "/testdata/schedulerdata/createactions"
	logger := logger.NewLogger("", true, true, true, true)
	connection := connection()
	scheduler := console.NewScheduler(filePath, []string{}, logger, connection, nil, context.Background())

	if err := scheduler.Run(); err != nil {
		t.Error(err)
	}

	exists, err := connection.Indexer(nil).Exists("my_index")

	if err != nil {
		t.Error(err)
	}

	assert.True(t, exists)

	if err := connection.Indexer(nil).DeleteIndex("my_index"); err != nil {
		t.Error(err)
	}

	exists, err = connection.Indexer(nil).Exists("my_async_index")

	if err != nil {
		t.Error(err)
	}

	assert.True(t, exists)

	if err := connection.Indexer(nil).DeleteIndex("my_async_index"); err != nil {
		t.Error(err)
	}
}
