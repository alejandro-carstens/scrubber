package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/rs/xid"
)

func TestWatchAction(t *testing.T) {
	for _, watchActionFilePath := range watchActionDataProvider() {
		if _, err := createTestIndex("/testdata/create_variants_index.yml"); err != nil {
			t.Error(err)
		}

		time.Sleep(time.Duration(int64(2)) * time.Second)

		connection := connection()

		if _, err := connection.Builder("variants-1992.06.02").Insert(makeVariants(1000)...); err != nil {
			t.Error(err)
		}

		time.Sleep(time.Duration(int64(1)) * time.Second)

		takeAction(watchActionFilePath, t)

		if err := connection.Indexer(nil).DeleteIndex("_all"); err != nil {
			t.Error(err)
		}
	}
}

func watchActionDataProvider() []string {
	return []string{
		"/testdata/watch_count.yml",
		"/testdata/watch_average_count.yml",
		"/testdata/watch_stats_count.yml",
		"/testdata/watch_count_no_min_threshold.yml",
	}
}

type variant struct {
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
	prices := []int{2000, 1001, 200, 1500, 5000}

	initialTime := time.Now()

	variants := []interface{}{}

	for i := 0; i < count; i++ {
		for index := range colors {
			initialTime = initialTime.Add(time.Duration(-1) * time.Minute)

			variant := &variant{
				Id:    xid.New().String(),
				Price: prices[index] - i,
			}

			variant.Attributes.Color = colors[index]
			variant.Attributes.Size = sizes[index]
			variant.Attributes.Sku = fmt.Sprintf("%v-%v-%v", colors[index], sizes[index], i)
			variant.CreatedAt = initialTime

			variants = append(variants, variant)
		}
	}

	return variants
}
