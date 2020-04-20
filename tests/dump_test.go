package tests

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/stretchr/testify/assert"
)

func TestDump(t *testing.T) {
	connection := connection()

	for _, data := range dumpDataProvider() {
		if _, err := createTestIndex("/testdata/create_variants_index.yml"); err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Duration(int64(2)) * time.Second)

		if _, err := connection.Builder("variants-1992.06.02").Insert(makeVariants(1000)...); err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Duration(int64(1)) * time.Second)

		takeAction(data["dump_file_path"].(string), t)

		builder := connection.Builder("variants-1992.06.02")

		builder.WhereNested("attributes.color", "=", "Red").
			FilterNested("attributes.size", "<=", 31).
			MatchInNested("attributes.sku", []interface{}{"Red-31"}).
			Where("price", "<", 150).
			Where("other_key", "<>", 300)

		count, err := builder.Count()

		if err != nil {
			t.Fatal(err)
		}

		verifyDataFiles(int(count), data["concurrency"].(int), t)

		verifyAliases(t)

		verifyMappings(t)

		verifySettings(t)

		if err := connection.Indexer(nil).DeleteIndex("_all"); err != nil {
			t.Fatal(err)
		}

		if err := os.RemoveAll("/usr/share/scrubber/scrubber_test"); err != nil {
			t.Fatal(err)
		}
	}
}

func dumpDataProvider() []map[string]interface{} {
	data := []map[string]interface{}{}

	data = append(data, map[string]interface{}{
		"dump_file_path": "/testdata/concurrent_dump.yml",
		"concurrency":    3,
	})
	data = append(data, map[string]interface{}{
		"dump_file_path": "/testdata/dump.yml",
		"concurrency":    1,
	})

	return data
}

func verifyDataFiles(expectedCount int, concurrency int, t *testing.T) {
	counter := 0

	for i := 0; i < concurrency; i++ {
		file, err := os.Open(fmt.Sprintf("/usr/share/scrubber/scrubber_test/variants-1992.06.02/data_%v.json", i))

		if err != nil {
			t.Fatal(err)
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			counter++
		}

		if err := scanner.Err(); err != nil {
			t.Fatal(err)
		}
	}

	assert.Equal(t, expectedCount, counter)
}

func verifyAliases(t *testing.T) {
	assertion, err := compareFiles(
		"/usr/share/scrubber/scrubber_test/variants-1992.06.02/aliases.json",
		"/go/src/scrubber/tests/testdata/dumpdata/aliases.json",
	)

	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, assertion)
}

func verifyMappings(t *testing.T) {
	assertion, err := compareFiles(
		"/usr/share/scrubber/scrubber_test/variants-1992.06.02/mappings.json",
		"/go/src/scrubber/tests/testdata/dumpdata/mappings.json",
	)

	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, assertion)
}

func verifySettings(t *testing.T) {
	c1, err := ioutil.ReadFile("/usr/share/scrubber/scrubber_test/variants-1992.06.02/settings.json")

	if err != nil {
		t.Fatal(err)
	}

	c2, err := ioutil.ReadFile("/go/src/scrubber/tests/testdata/dumpdata/settings.json")

	if err != nil {
		t.Fatal(err)
	}

	p1, err := gabs.ParseJSON([]byte(string(c1)))

	if err != nil {
		t.Fatal(err)
	}

	p2, err := gabs.ParseJSON([]byte(string(c2)))

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(
		t,
		p1.S("index", "number_of_replicas").Data().(string),
		p2.S("index", "number_of_replicas").Data().(string),
	)
	assert.Equal(
		t,
		p1.S("index", "number_of_shards").Data().(string),
		p2.S("index", "number_of_shards").Data().(string),
	)
	assert.Equal(
		t,
		p1.S("index", "provided_name").Data().(string),
		p2.S("index", "provided_name").Data().(string),
	)
}
