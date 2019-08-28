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

const SEED_INDEX_COUNT int = 4
const ELASTICSEARCH_BULK_INSERT_LIMIT int = 10000

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

	return actions.Create(context, logger.NewLogger("", true, true, true, true), connection())
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

func snapshotCleanup(repository, snapshot, index string, c *golastic.Connection) error {
	var conn *golastic.Connection

	if c == nil {
		conn = connection()
	} else {
		conn = c
	}

	if len(snapshot) > 0 {
		if _, err := conn.Indexer(nil).DeleteSnapshot(repository, snapshot); err != nil {
			return err
		}
	}

	if _, err := conn.Indexer(nil).DeleteRepositories(repository); err != nil {
		return err
	}

	return conn.Indexer(nil).DeleteIndex(index)
}

func takeActionAsync(path string, t *testing.T, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	takeAction(path, t)
}

func connection() *golastic.Connection {
	connection := golastic.NewConnection(&golastic.ConnectionContext{
		Urls:                []string{os.Getenv("ELASTICSEARCH_URI")},
		Password:            os.Getenv("ELASTICSEARCH_PASSWORD"),
		Username:            os.Getenv("ELASTICSEARCH_USERNAME"),
		HealthCheckInterval: 30,
	})

	if err := connection.Connect(); err != nil {
		panic(err)
	}

	return connection
}

func seedIndexAsync(index string, count int, connection *golastic.Connection, waitGroup *sync.WaitGroup, useConstantTime bool) {
	defer waitGroup.Done()

	inserts := []interface{}{}

	for i := 0; i < count; i++ {
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

		inserts = append(inserts, value)

		if count >= ELASTICSEARCH_BULK_INSERT_LIMIT && (i+1)%ELASTICSEARCH_BULK_INSERT_LIMIT == 0 {
			if _, err := connection.Builder(index).Insert(inserts...); err != nil {
				panic(err)
			}

			inserts = []interface{}{}
		}
	}

	if count < ELASTICSEARCH_BULK_INSERT_LIMIT {
		if _, err := connection.Builder(index).Insert(inserts...); err != nil {
			panic(err)
		}
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)
}
