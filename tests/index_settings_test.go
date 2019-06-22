package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/alejandro-carstens/golastic"
	"github.com/stretchr/testify/assert"
)

func TestIndexSettings(t *testing.T) {
	if _, err := createTestIndex("/../stubs/create_pattern_index.yml"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	takeAction("/../stubs/update_settings.yml", t)

	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		t.Error(err)
	}

	settings, err := builder.Settings("alejandro-carstens-1992.06.02")

	if err != nil {
		t.Error(err)
	}

	indexSettings, valid := settings["alejandro-carstens-1992.06.02"]

	if !valid {
		t.Error(errors.New("could not retrieve settings for index"))
	}

	assert.Equal(t, "true", indexSettings.S("index", "blocks", "write").Data().(string))
	assert.Equal(t, "7s", indexSettings.S("index", "refresh_interval").Data().(string))
	assert.Equal(t, "2", indexSettings.S("index", "number_of_replicas").Data().(string))

	if err := builder.DeleteIndex("alejandro-carstens-1992.06.02"); err != nil {
		t.Error(err)
	}
}
