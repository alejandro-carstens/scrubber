package tests

import (
	"sync"
	"testing"

	"github.com/alejandro-carstens/golastic"
)

func TestRollover(t *testing.T) {
	builder, err := golastic.NewBuilder(nil, nil)

	if err != nil {
		t.Error(err)
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	go seedIndexAsync("my_index-00001", 101, builder, &waitGroup, false)

	waitGroup.Wait()

	takeAction("/testdata/alias_rollover_index.yml", t)

	takeAction("/testdata/rollover.yml", t)

	if err := builder.DeleteIndex("my_index-00001"); err != nil {
		t.Error(err)
	}
}
