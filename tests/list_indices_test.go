package tests

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListIndices(t *testing.T) {
	var waitGroup sync.WaitGroup

	createIndexFiles := []string{
		"/testdata/create_index_by_name_1.yml",
		"/testdata/create_index_by_name_2.yml",
	}

	waitGroup.Add(len(createIndexFiles))

	for _, createIndexFile := range createIndexFiles {
		go createTestIndexAsync(createIndexFile, &waitGroup)
	}

	waitGroup.Wait()

	action := takeAction("/testdata/list_indices.yml", t)

	expectedIndices := []string{
		"my_index-2019-01-03",
		"my_index-2019-01-02",
	}

	assert.ElementsMatch(t, expectedIndices, action.List())

	if err := connection().Indexer(nil).DeleteIndex("_all"); err != nil {
		t.Error(err)
	}
}
