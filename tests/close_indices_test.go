package tests

import (
	"testing"
	"time"

	"github.com/alejandro-carstens/golastic"
	"github.com/stretchr/testify/assert"
)

func TestCloseIndices(t *testing.T) {
	if _, err := createTestIndex("/test_files/create_index.yml"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	takeAction("/test_files/close_indices_by_creation_date.yml", t)

	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		t.Error(err)
	}

	response, err := builder.IndexCat("my_index")

	if err != nil {
		t.Error(err)
	}

	children, err := response.Children()

	if err != nil {
		t.Error(err)
	}

	indexCat := children[0]

	assert.Equal(t, "close", indexCat.S("status").Data().(string))

	if err := builder.DeleteIndex("my_index"); err != nil {
		t.Error(err)
	}
}
