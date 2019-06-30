package tests

import (
	"sync"
	"testing"

	"github.com/alejandro-carstens/golastic"
	"github.com/stretchr/testify/assert"
)

func TestListIndices(t *testing.T) {
	var waitGroup sync.WaitGroup

	createIndexFiles := []string{
		"/test_files/create_index_by_name_1.yml",
		"/test_files/create_index_by_name_2.yml",
	}

	waitGroup.Add(len(createIndexFiles))

	for _, createIndexFile := range createIndexFiles {
		go createTestIndexAsync(createIndexFile, &waitGroup)
	}

	waitGroup.Wait()

	action := takeAction("/test_files/list_indices.yml", t)

	expectedIndices := []string{
		"my_index-2019-01-03",
		"my_index-2019-01-02",
	}

	assert.ElementsMatch(t, expectedIndices, action.List())

	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		t.Error(err)
	}

	for _, index := range expectedIndices {
		if err := builder.DeleteIndex(index); err != nil {
			t.Error(err)
		}
	}
}
