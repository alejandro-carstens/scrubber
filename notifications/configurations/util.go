package configurations

import "errors"

func Config(notificationType string) (Configurable, error) {
	var config Configurable

	switch notificationType {
	case "slack":
		config = &Slack{}
		break
	default:
		return nil, errors.New("invalid configuration type")
	}

	return config.FillFromEnvs().Validate()
}
