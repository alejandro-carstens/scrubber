package config

import (
	"errors"
	"strconv"
)

// GetFromEnvVars returns the configuration for the specified store
func GetFromEnvVars(store string) (Configurable, error) {
	var config Configurable

	switch store {
	case "mysql":
		config = &MySQL{}
		break
	default:
		return nil, errors.New("invalid store specified")
	}

	return config.FillFromEnvs().Validate()
}

func stringToInt(value string) int {
	v, _ := strconv.Atoi(value)

	return v
}
