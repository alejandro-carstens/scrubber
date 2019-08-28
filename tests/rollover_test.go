package tests

import (
	"sync"
	"testing"
)

func TestRollover(t *testing.T) {
	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	connection := connection()

	go seedIndexAsync("my_index-00001", 101, connection, &waitGroup, false)

	waitGroup.Wait()

	takeAction("/testdata/alias_rollover_index.yml", t)

	takeAction("/testdata/rollover.yml", t)

	if err := connection.Indexer(nil).DeleteIndex("my_index-00001"); err != nil {
		t.Error(err)
	}
}
