package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportDump(t *testing.T) {
	connection := connection()

	for _, data := range importDumpDataProvider() {
		takeAction(data["import_dump_file_path"].(string), t)

		// verify mappings
		// verify aliases
		// verify settings

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
	})

	return data
}
