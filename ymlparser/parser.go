package ymlparser

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Jeffail/gabs"
	"github.com/icza/dyno"
	"gopkg.in/yaml.v2"
)

func Parse(fileName string) (*gabs.Container, error) {
	reader, err := os.Open(fileName)

	defer reader.Close()

	if err != nil {
		return nil, err
	}

	fileContents, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	data := map[interface{}]interface{}{}

	if err := yaml.Unmarshal(fileContents, data); err != nil {
		return nil, err
	}

	encodedJSON, err := json.Marshal(dyno.ConvertMapI2MapS(data))

	if err != nil {
		return nil, err
	}

	container, err := gabs.ParseJSON(encodedJSON)

	if err != nil {
		return nil, err
	}

	return container, nil
}
