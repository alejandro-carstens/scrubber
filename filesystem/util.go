package filesystem

import "errors"

// Build returns a filesystem implementation
// based on the passed in configuration
func Build(config Configurable) (Storeable, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	switch config.Name() {
	case "local":
		return new(local).Init(config)
	}

	return nil, errors.New("invalid filesystem specified")
}
