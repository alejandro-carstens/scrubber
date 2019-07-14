package tests

import (
	"testing"

	"github.com/alejandro-carstens/golastic"
	"github.com/stretchr/testify/assert"
)

func TestCreateIndex(t *testing.T) {
	action, err := createTestIndex("/testfiles/create_index.yml")

	if err != nil {
		t.Error(err)
	}

	assert.False(t, action.HasErrors())

	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		t.Error(err)
	}

	if err := builder.DeleteIndex("my_index"); err != nil {
		t.Error(err)
	}
}
