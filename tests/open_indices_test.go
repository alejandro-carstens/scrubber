package tests

import (
	"testing"
	"time"

	"github.com/alejandro-carstens/golastic"
	"github.com/stretchr/testify/assert"
)

func TestOpenIndices(t *testing.T) {
	if _, err := createTestIndex("/../stubs/create_index.yml"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		t.Error(err)
	}

	if _, err := builder.Close("my_index"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	takeAction("/../stubs/open_closed_index.yml", t)

	response, err := builder.IndexCat("my_index")

	if err != nil {
		t.Error(err)
	}

	children, err := response.Children()

	if err != nil {
		t.Error(err)
	}

	indexCat := children[0]

	assert.Equal(t, "open", indexCat.S("status").Data().(string))

	if err := builder.DeleteIndex("my_index"); err != nil {
		t.Error(err)
	}
}
