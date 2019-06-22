package ymlparser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	currentPath, err := os.Getwd()

	if err != nil {
		t.Error(err)
	}

	data, err := Parse(currentPath + "/../stubs/parser_test.yml")

	if err != nil {
		t.Error(err)
	}

	response, err := data.ChildrenMap()

	if err != nil {
		t.Error(err)
	}

	expectedData := map[string]interface{}{
		"options": map[string]interface{}{
			"name": "my_index",
		},
		"extra_settings": map[string]interface{}{
			"settings": map[string]interface{}{
				"number_of_replicas": float64(3),
				"number_of_shards":   float64(2),
			},
		},
		"mappings": map[string]interface{}{
			"type1": map[string]interface{}{
				"properties": map[string]interface{}{
					"field1": map[string]interface{}{
						"index": "not_analyzed",
						"type":  "string",
					},
				},
			},
		},
	}

	assert.Equal(t, "create_index", response["action"].Data())
	assert.Equal(t, "Creates the specified index with the specified settings", response["description"].Data())
	assert.Equal(t, expectedData["options"], response["options"].Data())
	assert.Equal(t, expectedData["extra_settings"], response["extra_settings"].Data())
	assert.Equal(t, expectedData["mappings"], response["mappings"].Data())
}
