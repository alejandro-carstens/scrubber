package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCloseIndices(t *testing.T) {
	if _, err := createTestIndex("/testdata/create_index.yml"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	takeAction("/testdata/close_indices_by_creation_date.yml", t)

	connection := connection()

	response, err := connection.Indexer(nil).IndexCat("my_index")

	if err != nil {
		t.Error(err)
	}

	children, err := response.Children()

	if err != nil {
		t.Error(err)
	}

	indexCat := children[0]

	assert.Equal(t, "close", indexCat.S("status").Data().(string))

	if err := connection.Indexer(nil).DeleteIndex("my_index"); err != nil {
		t.Error(err)
	}
}
