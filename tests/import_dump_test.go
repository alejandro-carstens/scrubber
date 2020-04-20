package tests

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/alejandro-carstens/golastic"
	"github.com/stretchr/testify/assert"
)

func TestImportDump(t *testing.T) {
	connection := connection()

	for _, data := range importDumpDataProvider() {
		takeAction(data["import_dump_file_path"].(string), t)

		verifyImportDumpMappings(t, connection, data["index"].(string), data["mappings_file_path"].(string))
		verifyImportDumpAliases(t, connection, data["index"].(string), data["aliases_file_path"].(string))
		verifyImportDumpSettings(t, connection, data["index"].(string), data["expected_settings"].(map[string]string))

		builder := connection.Builder(data["index"].(string))

		builder.WhereNested("attributes.color", "=", "Red").
			FilterNested("attributes.size", "<=", 31).
			MatchInNested("attributes.sku", []interface{}{"Red-31"}).
			Where("price", "<", 150).
			Where("other_key", "<>", 300)

		count, err := builder.Count()

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, data["expected_count"].(int64), count)

		if err := connection.Indexer(nil).DeleteIndex("_all"); err != nil {
			t.Fatal(err)
		}
	}
}

func importDumpDataProvider() []map[string]interface{} {
	data := []map[string]interface{}{}

	data = append(data, map[string]interface{}{
		"import_dump_file_path": "/testdata/import_dump.yml",
		"index":                 "import-scrubber_test-variants-1992.06.02",
		"expected_count":        int64(148),
		"mappings_file_path":    "/go/src/scrubber/tests/testdata/importdumpdata/responses/mappings.json",
		"aliases_file_path":     "/go/src/scrubber/tests/testdata/importdumpdata/responses/aliases.json",
		"expected_settings": map[string]string{
			"number_of_replicas": "0",
			"number_of_shards":   "2",
			"provided_name":      "import-scrubber_test-variants-1992.06.02",
		},
	})

	return data
}

func verifyImportDumpMappings(t *testing.T, connection *golastic.Connection, index, mappingsFilePath string) {
	response, err := connection.Indexer(nil).Mappings(index)

	if err != nil {
		t.Fatal(err)
	}

	contents, err := ioutil.ReadFile(mappingsFilePath)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(contents), response.String())
}

func verifyImportDumpAliases(t *testing.T, connection *golastic.Connection, index, aliasesFilePath string) {
	response, err := connection.Indexer(nil).Aliases(index)

	if err != nil {
		t.Fatal(err)
	}

	contents, err := ioutil.ReadFile(aliasesFilePath)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(contents), response.String())
}

func verifyImportDumpSettings(t *testing.T, connection *golastic.Connection, index string, expectedSettings map[string]string) {
	response, err := connection.Indexer(nil).Settings(index)

	if err != nil {
		t.Fatal(err)
	}

	settings, valid := response[index]

	if !valid {
		t.Fatal(fmt.Sprintf("no settings found for index: %v", index))
	}

	assert.Equal(t, expectedSettings["number_of_replicas"], settings.S("index", "number_of_replicas").Data().(string))
	assert.Equal(t, expectedSettings["number_of_shards"], settings.S("index", "number_of_shards").Data().(string))
	assert.Equal(t, expectedSettings["provided_name"], settings.S("index", "provided_name").Data().(string))
}
