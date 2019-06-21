package tests

import (
	"log"
	"os"
	"scrubber/actions"
	"scrubber/actions/contexts"
	"scrubber/logging"
	yml "scrubber/yml_parser"
	"sync"
	"testing"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/golastic"
	"github.com/stretchr/testify/assert"
)

func createTestIndex(filePath string) (actions.Actionable, error) {
	config, err := config(filePath)

	if err != nil {
		return nil, err
	}

	action, err := getAction(config)

	if err != nil {
		return nil, err
	}

	return action.Perform(), nil
}

func createTestIndexAsync(filePath string, waitGroup *sync.WaitGroup) {
	log.Println("Creating " + filePath)

	if _, err := createTestIndex(filePath); err != nil {
		log.Println(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	defer waitGroup.Done()
}

func config(path string) (*gabs.Container, error) {
	currentPath, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	return yml.Parse(currentPath + path)
}

func getAction(config *gabs.Container) (actions.Actionable, error) {
	context, err := contexts.New(config)

	if err != nil {
		return nil, err
	}

	return actions.Create(context, logging.NewSrvLogger("", true, true, true, true))
}

func takeAction(path string, t *testing.T) {
	config, err := config(path)

	if err != nil {
		t.Error(err)
	}

	action, err := getAction(config)

	if err != nil {
		t.Error(err)
	}

	assert.False(t, action.Perform().HasErrors())

	time.Sleep(time.Duration(int64(2)) * time.Second)
}

func snapshotCleanup(repository, snapshot, index string, builder *golastic.ElasticsearchBuilder) error {
	var err error

	if builder == nil {
		builder, err = golastic.NewBuilder(nil, nil)

		if err != nil {
			return err
		}
	}

	if len(snapshot) > 0 {
		if _, err := builder.DeleteSnapshot(repository, snapshot); err != nil {
			return err
		}
	}

	if _, err := builder.DeleteRepositories(repository); err != nil {
		return err
	}

	return builder.DeleteIndex(index)
}
