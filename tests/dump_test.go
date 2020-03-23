package tests

import (
	"testing"
	"time"
)

func TestDump(t *testing.T) {
	connection := connection()

	if _, err := createTestIndex("/testdata/create_variants_index.yml"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	if _, err := connection.Builder("variants-1992.06.02").Insert(makeVariants(1000)...); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(1)) * time.Second)

	takeAction("/testdata/dump.yml", t)

	if err := connection.Indexer(nil).DeleteIndex("_all"); err != nil {
		t.Error(err)
	}
}
