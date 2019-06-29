package tests

import (
	"testing"
	"time"

	"github.com/alejandro-carstens/golastic"
	"github.com/stretchr/testify/assert"
)

func TestAlias(t *testing.T) {
	if _, err := createTestIndex("/test_files/create_pattern_index.yml"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	takeAction("/test_files/add_alias.yml", t)

	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		t.Error(err)
	}

	response, err := builder.AliasesCat()

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1, len(response))
	assert.Equal(t, "alejandro-carstens-1992.06.02", response[0].Index)
	assert.Equal(t, "*", response[0].Filter)
	assert.Equal(t, "great_alias", response[0].Alias)
	assert.Equal(t, "1", response[0].RoutingIndex)
	assert.Equal(t, "1,2,3", response[0].RoutingSearch)

	takeAction("/test_files/remove_alias.yml", t)

	response, err = builder.AliasesCat()

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 0, len(response))

	if err := builder.DeleteIndex("alejandro-carstens-1992.06.02"); err != nil {
		t.Error(err)
	}
}
