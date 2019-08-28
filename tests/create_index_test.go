package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndex(t *testing.T) {
	action, err := createTestIndex("/testdata/create_index.yml")

	if err != nil {
		t.Error(err)
	}

	assert.False(t, action.HasErrors())

	if err := connection().Indexer(nil).DeleteIndex("my_index"); err != nil {
		t.Error(err)
	}
}
