package tests

import (
	"log"
	"os"
	"scrubber/actions"
	"scrubber/actions/contexts"
	"scrubber/logger"
	"scrubber/ymlparser"
	"strconv"
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

	return ymlparser.Parse(currentPath + path)
}

func getAction(config *gabs.Container) (actions.Actionable, error) {
	context, err := contexts.New(config)

	if err != nil {
		return nil, err
	}

	return actions.Create(context, logger.NewLogger("", true, true, true, true), nil)
}

func takeAction(path string, t *testing.T) actions.Actionable {
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

	return action
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

func takeActionAsync(path string, t *testing.T, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	takeAction(path, t)
}

func seedIndexAsync(index string, count int, builder *golastic.ElasticsearchBuilder, waitGroup *sync.WaitGroup, useConstantTime bool) {
	defer waitGroup.Done()

	inserts := []golastic.ElasticModelable{}

	for i := 0; i < count; i++ {
		insert := golastic.NewGolasticModel()

		value := map[string]interface{}{}

		value["id"] = strconv.Itoa(i + 1 + 1000000000)
		value["exception"] = "Exception exception exception exception exception exception exception exception exception"
		value["request"] = "Request request request request request request request request request"
		value["message"] = "Message message message message message message message message message"
		value["bytes"] = int64(i)
		value["number"] = float64(i)

		if useConstantTime {
			constantTime, err := time.Parse(time.RFC3339, "2017-11-12T11:45:26.371Z")

			if err != nil {
				panic(err)
			}

			value["created_at"] = constantTime
		} else {
			value["created_at"] = time.Now().Add(time.Duration(int64(-1*(i+1))) * time.Hour)
		}

		insert.SetData(value)
		insert.SetIndex(index)

		inserts = append(inserts, insert)

		if count >= ELASTICSEARCH_BULK_INSERT_LIMIT && (i+1)%ELASTICSEARCH_BULK_INSERT_LIMIT == 0 {
			if _, err := builder.Insert(inserts...); err != nil {
				panic(err)
			}

			inserts = []golastic.ElasticModelable{}
		}
	}

	if count < ELASTICSEARCH_BULK_INSERT_LIMIT {
		if _, err := builder.Insert(inserts...); err != nil {
			panic(err)
		}
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)
}
