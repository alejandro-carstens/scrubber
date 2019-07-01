package options

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
)

func toContainer(val interface{}) *gabs.Container {
	b, _ := json.Marshal(val)

	container, _ := gabs.ParseJSON(b)

	return container
}
