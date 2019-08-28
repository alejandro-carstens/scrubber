package tests

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAlias(t *testing.T) {
	if _, err := createTestIndex("/testdata/create_pattern_index.yml"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	takeAction("/testdata/add_alias.yml", t)

	connection := connection()

	response, err := connection.Indexer(nil).AliasesCat()

	if err != nil {
		t.Error(err)
	}

	aliases, err := response.Children()

	if err != nil {
		t.Error(err)
	}

	log.Println(aliases[0].String())

	assert.Equal(t, 1, len(aliases))
	assert.Equal(t, "alejandro-carstens-1992.06.02", aliases[0].S("index").Data().(string))
	assert.Equal(t, "*", aliases[0].S("filter").Data().(string))
	assert.Equal(t, "great_alias", aliases[0].S("alias").Data().(string))
	assert.Equal(t, "1", aliases[0].S("routing.index").Data().(string))
	assert.Equal(t, "1,2,3", aliases[0].S("routing.search").Data().(string))

	takeAction("/testdata/remove_alias.yml", t)

	response, err = connection.Indexer(nil).AliasesCat()

	if err != nil {
		t.Error(err)
	}

	aliases, err = response.Children()

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 0, len(aliases))

	if err := connection.Indexer(nil).DeleteIndex("alejandro-carstens-1992.06.02"); err != nil {
		t.Error(err)
	}
}
