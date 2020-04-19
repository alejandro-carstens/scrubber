package tests

import (
	"testing"
)

func TestImportDump(t *testing.T) {
	connection := connection()

	for _, data := range importDumpDataProvider() {
		takeAction(data["import_dump_file_path"].(string), t)

		// builder := connection.Builder(data["index"].(string))

		// builder.WhereNested("attributes.color", "=", "Red").
		// 	FilterNested("attributes.size", "<=", 31).
		// 	MatchInNested("attributes.sku", []interface{}{"Red-31"}).
		// 	Where("price", "<", 150).
		// 	Where("other_key", "<>", 300)

		// count, err := builder.Count()

		// if err != nil {
		// 	t.Fatal(err)
		// }

		if err := connection.Indexer(nil).DeleteIndex("_all"); err != nil {
			t.Fatal(err)
		}
	}
}

func importDumpDataProvider() []map[string]interface{} {
	data := []map[string]interface{}{}

	data = append(data, map[string]interface{}{
		"import_dump_file_path": "/testdata/import_dump.yml",
	})

	return data
}
