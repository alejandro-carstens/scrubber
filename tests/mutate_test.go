package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMutateUpdateAction(t *testing.T) {
	connection := connection()

	if _, err := createTestIndex("/testdata/create_variants_index.yml"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	if _, err := connection.Builder("variants-1992.06.02").Insert(makeVariants(1000)...); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(1)) * time.Second)

	takeAction("/testdata/mutate_update.yml", t)

	time.Sleep(time.Duration(int64(1)) * time.Second)

	builder := connection.Builder("variants-1992.06.02")

	builder.Where("key", "=", "value").Where("other_key", "=", 300)

	count, err := builder.Count()

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, int64(148), count)

	if err := connection.Indexer(nil).DeleteIndex("_all"); err != nil {
		t.Error(err)
	}
}

func TestMutateDeleteAction(t *testing.T) {
	connection := connection()

	if _, err := createTestIndex("/testdata/create_variants_index.yml"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	if _, err := connection.Builder("variants-1992.06.02").Insert(makeVariants(1000)...); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(1)) * time.Second)

	takeAction("/testdata/mutate_delete.yml", t)

	time.Sleep(time.Duration(int64(1)) * time.Second)

	builder := connection.Builder("variants-1992.06.02")

	builder.WhereNested("attributes.color", "=", "Red").
		FilterNested("attributes.size", "<=", 31).
		MatchInNested("attributes.sku", []interface{}{"Red-31"}).
		Where("price", "<", 150)

	count, err := builder.Count()

	assert.Nil(t, err)
	assert.Equal(t, int64(0), count)

	if err := connection.Indexer(nil).DeleteIndex("_all"); err != nil {
		t.Error(err)
	}
}
