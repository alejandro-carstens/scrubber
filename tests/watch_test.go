package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/rs/xid"
)

func TestWatchCountAction(t *testing.T) {
	if _, err := createTestIndex("/testdata/create_variants_index.yml"); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(2)) * time.Second)

	connection := connection()

	if _, err := connection.Builder("variants-1992.06.02").Insert(makeVariants(1000)...); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(int64(1)) * time.Second)

	takeAction("/testdata/watch_count.yml", t)

	if err := connection.Indexer(nil).DeleteIndex("variants-1992.06.02"); err != nil {
		t.Error(err)
	}
}

type Variant struct {
	Id         string    `json:"id"`
	Price      int       `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
	Attributes struct {
		Color string `json:"color"`
		Size  int    `json:"size"`
		Sku   string `json:"sku"`
	} `json:"attributes"`
}

func makeVariants(count int) []interface{} {
	colors := []string{"Blue", "Red", "Red", "Purple", "Black"}
	sizes := []int{30, 31, 32, 33, 34}
	prices := []int{200, 150, 125, 100, 85}

	initialTime := time.Now()

	variants := []interface{}{}

	for i := 0; i < count; i++ {
		for index := range colors {
			variant := &Variant{
				Id:    xid.New().String(),
				Price: prices[index] - i,
			}

			variant.Attributes.Color = colors[index]
			variant.Attributes.Size = sizes[index]
			variant.Attributes.Sku = fmt.Sprintf("%v-%v-%v", colors[index], sizes[index], i)
			variant.CreatedAt = initialTime.Add(time.Duration(1) * time.Hour)

			variants = append(variants, variant)
		}
	}

	return variants
}
