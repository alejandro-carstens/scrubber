package options

import (
	"encoding/json"

	"github.com/spf13/pflag"

	"github.com/Jeffail/gabs"
)

func toContainer(val interface{}) *gabs.Container {
	b, _ := json.Marshal(val)

	container, _ := gabs.ParseJSON(b)

	return container
}

func stringFromFlags(flags *pflag.FlagSet, key string) string {
	value, _ := flags.GetString(key)

	return value
}

func boolFromFlags(flags *pflag.FlagSet, key string) bool {
	value, _ := flags.GetBool(key)

	return value
}

func intFromFlags(flags *pflag.FlagSet, key string) int {
	value, _ := flags.GetInt(key)

	return value
}
